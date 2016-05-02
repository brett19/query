//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"runtime"

	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/util"
	"github.com/couchbase/query/value"
	"fmt"
)

type Parallel struct {
	base
	plan         *plan.Parallel
	child        Operator
	inChannel    chan value.AnnotatedValue
	outChannel   chan value.AnnotatedValue
}

func NewParallel(plan *plan.Parallel, child Operator) *Parallel {
	rv := &Parallel{
		base:         newBase(),
		plan:         plan,
		child:        child,
		inChannel:    make(chan value.AnnotatedValue, 40),
		outChannel:   make(chan value.AnnotatedValue, 40),
	}

	rv.output = rv
	return rv
}

func (this *Parallel) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitParallel(this)
}

func (this *Parallel) Copy() Operator {
	return &Parallel{
		base:         this.base.copy(),
		plan:         this.plan,
		child:        this.child.Copy(),
	}
}

func (this *Parallel) RunOnce(context *Context, parent value.Value) {
	n := util.MinInt(this.plan.MaxParallelism(), context.MaxParallelism())
	children := make([]*parallelChild, n)

	doneCh := make(chan bool)

	children[0] = newParallelChild(this, this.outChannel)
	for i := 1; i < n; i++ {
		children[i] = children[0].Copy().(*parallelChild)
	}

	for i := 0; i < n; i++ {
		i := i
		go func() {
			children[i].Start(context, parent)
			doneCh <- true
		}()
	}

	go func() {
		this.input.RunOnce(context, parent)
		close(this.inChannel)

		for i := 0; i < n; i++ {
			<- doneCh
		}
		close(this.outChannel)
	}()

	Loop:
	for {
		item, ok := <-this.outChannel
		if !ok {
			break Loop
		}
		this.sendItem(item, context)
	}
}

func (this *Parallel) Item(item value.AnnotatedValue, context *Context) bool {
	this.inChannel <- item
	return true
}

type parallelChild struct {
	base
	parent        *Parallel
	child        Operator
}

func newParallelChild(parent *Parallel, itemCh chan value.AnnotatedValue) *parallelChild {
	rv := &parallelChild{
		base: newBase(),
		parent: parent,
		child: parent.child.Copy(),
	}

	rv.output = parent.output
	return rv
}

func (this *parallelChild) Accept(visitor Visitor) (interface{}, error) {
	panic(fmt.Sprintf("Internal operator parallelChild visited by %v.", visitor))
}

func (this *parallelChild) Copy() Operator {
	return &parallelChild{
		base: this.base.copy(),
		parent: this.parent,
		child: this.child.Copy(),
	}
}

func (this *parallelChild) Start(context *Context, parent value.Value) {
	this.child.SetInput(this)
	this.child.SetOutput(this)
	this.child.RunOnce(context, parent)
}

func (this *parallelChild) Item(item value.AnnotatedValue, context *Context) bool {
	this.parent.outChannel <- item
	return true
}

func (this *parallelChild) RunOnce(context *Context, parent value.Value) {
	// Need to pull items off and send them...
	for {
		item, ok := <- this.parent.inChannel
		if !ok {
			break
		}
		this.child.Item(item, context)
	}
}

var _PARALLEL_POOL = NewOperatorPool(runtime.NumCPU())
