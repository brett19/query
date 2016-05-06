//  Copieright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package value

import (
	"bytes"
)

type AnnotatedChannel chan AnnotatedValue
type AnnotatedValues []AnnotatedValue

func (this AnnotatedValues) Append(val AnnotatedValue, pool *AnnotatedPool) AnnotatedValues {
	if len(this) == cap(this) {
		avs := make(AnnotatedValues, len(this), len(this)<<1)
		copy(avs, this)
		pool.Put(this)
		this = avs
	}

	this = append(this, val)
	return this
}

/*
AnnotatedValue is a Value that can hold attachments and can hold data
from covering indexes.
*/
type AnnotatedValue interface {
	Value
	GetValue() Value
	Attachments() map[string]interface{}
	GetAttachment(key string) interface{}
	SetAttachment(key string, val interface{})
	Covers() map[string]Value
	GetCover(key string) Value
	SetCover(key string, val Value)
}

func NewAnnotatedValue(val interface{}) AnnotatedValue {
	switch val := val.(type) {
	case AnnotatedValue:
		return val
	case Value:
		av := &annotatedValue{
			Value: val,
		}
		return av
	default:
		return NewAnnotatedValue(NewValue(val))
	}
}

type annotatedValue struct {
	Value
	attachments map[string]interface{}
	covers      map[string]Value
}

func (this *annotatedValue) String() string {
	return this.Value.String()
}

func (this *annotatedValue) MarshalJSON() ([]byte, error) {
	return this.Value.MarshalJSON()
}

func (this *annotatedValue) FastMarshalJSON(buf *bytes.Buffer) error {
	return this.Value.FastMarshalJSON(buf)
}

func (this *annotatedValue) Copy() Value {
	return &annotatedValue{
		Value:       this.Value.Copy(),
		attachments: copyMap(this.attachments, self),
		covers:      this.covers,
	}
}

func (this *annotatedValue) CopyForUpdate() Value {
	return &annotatedValue{
		Value:       this.Value.CopyForUpdate(),
		attachments: copyMap(this.attachments, self),
		covers:      this.covers,
	}
}

func (this *annotatedValue) GetValue() Value {
	return this.Value
}

func (this *annotatedValue) Attachments() map[string]interface{} {
	return this.attachments
}

func (this *annotatedValue) GetAttachment(key string) interface{} {
	if this.attachments != nil {
		return this.attachments[key]
	}

	return nil
}

func (this *annotatedValue) SetAttachment(key string, val interface{}) {
	if this.attachments == nil {
		this.attachments = make(map[string]interface{})
	}

	this.attachments[key] = val
}

func (this *annotatedValue) Covers() map[string]Value {
	return this.covers
}

func (this *annotatedValue) GetCover(key string) Value {
	if this.covers != nil {
		return this.covers[key]
	}

	return nil
}

func (this *annotatedValue) SetCover(key string, val Value) {
	if this.covers == nil {
		this.covers = make(map[string]Value)
	}

	this.covers[key] = val
}
