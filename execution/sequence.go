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
	"github.com/couchbase/query/value"
)

type Sequence struct {
	base
	children     []Operator
}

func NewSequence(children ...Operator) *Sequence {
	rv := &Sequence{
		base:         newBase(),
		children:     children,
	}

	rv.output = rv
	return rv
}

func (this *Sequence) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSequence(this)
}

func (this *Sequence) Copy() Operator {
	children := _SEQUENCE_POOL.Get()

	for _, child := range this.children {
		children = append(children, child.Copy())
	}

	return &Sequence{
		base:         this.base.copy(),
		children:     children,
	}
}

func (this *Sequence) RunOnce(context *Context, parent value.Value) {
	n := len(this.children)

	input := this.input

	for i := 0; i < n; i++ {
		curr := this.children[i]
		curr.SetInput(input)
		if input != nil {
			input.SetOutput(curr)
		}
		input = curr
	}

	input.SetOutput(this.output)
	input.RunOnce(context, parent)
}

func (this *Sequence) Item(item value.AnnotatedValue, context *Context) bool {
	return this.children[0].Item(item, context)
}

var _SEQUENCE_POOL = NewOperatorPool(32)
