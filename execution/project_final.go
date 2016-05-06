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

type FinalProject struct {
	base
}

func NewFinalProject() *FinalProject {
	rv := &FinalProject{
		base: newBase(),
	}

	rv.output = rv
	return rv
}

func (this *FinalProject) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFinalProject(this)
}

func (this *FinalProject) Copy() Operator {
	return &FinalProject{this.base.copy()}
}

func (this *FinalProject) RunOnce(context *Context, parent value.Value) {
	this.runConsumer(this, context, parent)
}

func (this *FinalProject) Item(item value.AnnotatedValue, context *Context) bool {
	pv := item.GetAttachment("projection")
	if pv != nil {
		item.SetAttachment("projection", nil)
		_VALUEPOOL.Release(item)

		v := pv.(value.Value)
		return this.sendItem(value.NewAnnotatedValue(v), context)
	}

	return this.sendItem(item, context)
}
