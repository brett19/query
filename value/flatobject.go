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
	"github.com/couchbase/query/ssjson"
)

type flatPair struct {
	key   string
	value interface{}
}

type flatMap struct {
	items []flatPair
}

func (this *flatMap) Clear() {
	this.items = this.items[:0]
}

func (this *flatMap) Get(key string) interface{} {
	for _, k := range this.items {
		if k.key == key {
			return k.value
		}
	}
	return nil
}

func (this *flatMap) Set(key string, val interface{}) {
	for _, k := range this.items {
		if k.key == key {
			k.value = val
			return
		}
	}
	this.items = append(this.items, flatPair{
		key: key,
		value: val,
	})
}

func (this *flatMap) Unset(key string) {
	for i, v := range this.items {
		if v.key == key {
			copy(this.items[i:], this.items[i+1:])
			return
		}
	}
}

func (this *flatMap) Items() []flatPair {
	return this.items
}

func NewFlatObject(len int) Value {
	return &flatObject{}
}

type flatObject struct {
	items      flatMap
	parsed     Value
}

func (this *flatObject) String() string {
	return this.unwrap().String()
}

func (this *flatObject) MarshalJSON() ([]byte, error) {
	if this.parsed != nil {
		return this.parsed.MarshalJSON()
	}
	var buf bytes.Buffer
	err := this.FastMarshalJSON(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *flatObject) FastMarshalJSON(buf *bytes.Buffer) error {
	if this.parsed != nil {
		return this.parsed.FastMarshalJSON(buf)
	}
	buf.Write([]byte("{"))
	for i, v := range this.items.Items() {
		if i > 0 {
			buf.Write([]byte(","))
		}
		json.FastMarshal(buf, v.key)
		buf.Write([]byte(":"))
		json.FastMarshal(buf, v.value)
	}
	buf.Write([]byte("}"))
	return nil
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
	if this.parsed != nil {
		return this.parsed.Field(field)
	}
	for _, v := range this.items.Items() {
		if v.key == field {
			return v.value.(Value), true
		}
	}
	return missingField(field), false
}

func (this *flatObject) SetField(field string, val interface{}) error {
	if this.parsed != nil {
		return this.parsed.SetField(field, val)
	}
	this.items.Set(field, val)
	return nil
}

func (this *flatObject) UnsetField(field string) error {
	if this.parsed != nil {
		return this.parsed.UnsetField(field)
	}
	this.items.Unset(field)
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
	if this.parsed != nil {
		return this.parsed.FieldNames(buffer)
	}
	for i, k := range this.items.Items() {
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
		parsed := make(map[string]interface{}, len(this.items.Items()))
		for _, k := range this.items.Items() {
			parsed[k.key] = k.value
		}
		this.parsed = NewValue(parsed)
	}
	return this.parsed
}
