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
	"github.com/couchbase/query/util"
)

func NewFlatObject(len int) Value {
	return &flatObject{
		items: make([]flatPair, 0, len),
		parsed: nil,
	}
}

type flatPair struct {
	key   string
	value Value
}

type flatObject struct {
	items      []flatPair
	parsed     Value
}

func (this *flatObject) String() string {
	return this.unwrap().String()
}

func (this *flatObject) MarshalJSON() ([]byte, error) {
	return this.unwrap().MarshalJSON()
}

func (this *flatObject) FastMarshalJSON(buf *bytes.Buffer) error {
	return this.unwrap().FastMarshalJSON(buf)
}

func (this *flatObject) Type() Type {
	return OBJECT
}

func (this *flatObject) Actual() interface{} {
	return this.unwrap().Actual()
}

func (this *flatObject) Equals(other Value) Value {
	return this.unwrap().Equals(other)
}

func (this *flatObject) Collate(other Value) int {
	return this.unwrap().Collate(other)
}

func (this *flatObject) Compare(other Value) Value {
	return this.unwrap().Compare(other)
}

func (this *flatObject) Truth() bool {
	return this.unwrap().Truth()
}

func (this *flatObject) Copy() Value {
	return this.unwrap().Copy()
}

func (this *flatObject) CopyForUpdate() Value {
	return this.unwrap().CopyForUpdate()
}

func (this *flatObject) Field(field string) (Value, bool) {
	for _, v := range this.items {
		if v.key == field {
			return v.value, true
		}
	}
	return missingField(field), false
}

func (this *flatObject) SetField(field string, val interface{}) error {
	for _, v := range this.items {
		if v.key == field {
			v.value = NewValue(val)
			return nil
		}
	}
	this.items = append(this.items, flatPair{field, NewValue(val)})
	return nil
}

func (this *flatObject) UnsetField(field string) error {
	for i, v := range this.items {
		if v.key == field {
			copy(this.items[i:], this.items[i+1:])
			return nil
		}
	}
	return nil
}

func (this *flatObject) Index(index int) (Value, bool) {
	return missingIndex(index), false
}

func (this *flatObject) SetIndex(index int, val interface{}) error {
	return Unsettable(index)
}

func (this *flatObject) Slice(start, end int) (Value, bool) {
	return NULL_VALUE, false
}

func (this *flatObject) SliceTail(start int) (Value, bool) {
	return NULL_VALUE, false
}

func (this *flatObject) Descendants(buffer []interface{}) []interface{} {
	return this.unwrap().Descendants(buffer)
}

func (this *flatObject) Fields() map[string]interface{} {
	return this.unwrap().Fields()
}

func (this *flatObject) FieldNames(buffer []string) []string {
	for i, k := range this.items {
		buffer[i] = k.key
	}
	return buffer
}

func (this *flatObject) DescendantPairs(buffer []util.IPair) []util.IPair {
	return this.unwrap().DescendantPairs(buffer)
}

func (this *flatObject) Successor() Value {
	return this.unwrap().Successor()
}

func (this *flatObject) unwrap() Value {
	if this.parsed == nil {
		parsed := make(map[string]interface{}, len(this.items))
		for _, k := range this.items {
			parsed[k.key] = k.value
		}
		this.parsed = NewValue(parsed)
	}
	return this.parsed
}
