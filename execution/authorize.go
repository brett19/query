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
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

type Authorize struct {
	base
	plan         *plan.Authorize
	child        Operator
}

func NewAuthorize(plan *plan.Authorize, child Operator) *Authorize {
	rv := &Authorize{
		base:         newBase(),
		plan:         plan,
		child:        child,
	}

	rv.output = rv
	return rv
}

func (this *Authorize) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAuthorize(this)
}

func (this *Authorize) Copy() Operator {
	return &Authorize{
		base:         this.base.copy(),
		plan:         this.plan,
		child:        this.child.Copy(),
	}
}

func (this *Authorize) RunOnce(context *Context, parent value.Value) {
	this.child.SetInput(this.input)
	this.child.SetOutput(this)
	this.child.RunOnce(context, parent)
}

func (this *Authorize) Item(item value.AnnotatedValue, context *Context) bool {
	return this.output.Item(item, context)
}
