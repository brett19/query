//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execute

import (
	"fmt"

	"github.com/couchbaselabs/query/err"
	"github.com/couchbaselabs/query/value"
)

type Stream struct {
	base
}

func NewStream() *Stream {
	rv := &Stream{
		base: newBase(),
	}

	rv.output = rv
	return rv
}

func (this *Stream) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitStream(this)
}

func (this *Stream) Copy() Operator {
	return &Stream{this.base.copy()}
}

func (this *Stream) RunOnce(context *Context, parent value.Value) {
	this.runConsumer(this, context, parent)
}

func (this *Stream) processItem(item value.AnnotatedValue, context *Context) bool {
	mv := item.GetAttachment("meta")
	if mv == nil {
		context.ErrorChannel() <- err.NewError(nil, "Unable to find meta.")
		return false
	}

	meta := mv.(map[string]interface{})
	project, ok := meta["project"]
	if !ok {
		context.ErrorChannel() <- err.NewError(nil, "Unable to find projection.")
		return false
	}

	switch project := project.(type) {
	case value.Value:
		if project.Type() != value.MISSING {
			return context.Stream(project)
		} else {
			return true
		}
	default:
		context.ErrorChannel() <- err.NewError(nil,
			fmt.Sprintf("Unable to project value %v of type %T.", project, project))
		return false
	}
}
