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
	"encoding/json"
)

/*
sliceValue is defined as a slice of interfaces.
*/
type sliceValue []interface{}

/*
EMPTY_ARRAY_VALUE is initialized as a slice of interface.
*/
var EMPTY_ARRAY_VALUE = NewValue([]interface{}{})

func (this sliceValue) MarshalJSON() ([]byte, error) {
	return marshalArray(this)
}

/*
Type ARRAY
*/
func (this sliceValue) Type() Type {
	return ARRAY
}

/*
Cast receiver to an interface and return it.
*/
func (this sliceValue) Actual() interface{} {
	return ([]interface{})(this)
}

func (this sliceValue) Equals(other Value) Value {
	other = other.unwrap()
	switch other := other.(type) {
	case missingValue:
		return other
	case *nullValue:
		return other
	case sliceValue:
		return arrayEquals(this, other)
	case *listValue:
		return arrayEquals(this, other.slice)
	}

	return FALSE_VALUE
}

func (this sliceValue) Collate(other Value) int {
	other = other.unwrap()
	switch other := other.(type) {
	case sliceValue:
		return arrayCollate(this, other)
	case *listValue:
		return arrayCollate(this, other.slice)
	default:
		return int(ARRAY - other.Type())
	}
}

func (this sliceValue) Compare(other Value) Value {
	other = other.unwrap()
	switch other := other.(type) {
	case missingValue:
		return other
	case *nullValue:
		return other
	case sliceValue:
		return arrayCompare(this, other)
	case *listValue:
		return arrayCompare(this, other.slice)
	default:
		return NewValue(int(ARRAY - other.Type()))
	}
}

/*
If length of the slice  greater than 0, its a valid slice.
Return true.
*/
func (this sliceValue) Truth() bool {
	return len(this) > 0
}

/*
Call copySlice on the receiver and self and cast it to a
sliceValue.
*/
func (this sliceValue) Copy() Value {
	return sliceValue(copySlice(this, self))
}

/*
Call copySlice on the receiver and copyForUpdate, return a
pointer to a list value encapsulating it.This allows for a
copy for every element of the array by calling its
CopyForUpdate function.
*/
func (this sliceValue) CopyForUpdate() Value {
	return &listValue{copySlice(this, copyForUpdate)}
}

/*
Calls missingField.
*/
func (this sliceValue) Field(field string) (Value, bool) {
	return missingField(field), false
}

/*
Not valid for array/slice.
*/
func (this sliceValue) SetField(field string, val interface{}) error {
	return Unsettable(field)
}

/*
Not valid for array/slice.
*/
func (this sliceValue) UnsetField(field string) error {
	return Unsettable(field)
}

/*
If the input index is negative then count the index from
the last element.
*/
func (this sliceValue) Index(index int) (Value, bool) {
	if index < 0 {
		index = len(this) + index
	}

	if index >= 0 && index < len(this) {
		return NewValue(this[index]), true
	}

	return missingIndex(index), false
}

/*
If index is negative, add to the length and get the actual
index. In the event the new adjusted index is less than 0
or greater than/equal to the length of the slice return
Unsettable since Slices do NOT extend beyond length. If it
is a valid index, check the type of value. If it is a
missing value, set it to nil (do not add this field) and
if anything else, add the value at the particular index.
For all other cases, return a nil.
*/
func (this sliceValue) SetIndex(index int, val interface{}) error {
	if index < 0 {
		index = len(this) + index
	}

	if index < 0 || index >= len(this) {
		return Unsettable(index)
	}

	switch val := val.(type) {
	case missingValue:
		this[index] = nil
	default:
		this[index] = val
	}

	return nil
}

/*
If the start and/or end index is -ve, as per the N1QL specs,
add it to the length to get the actual index (from the end).
If it is a valid slice (start<=end, start >=0 and end less
than the length), return the slice by creating a valid value
and also return true. If the indices are not valid return a
missing value and false.
*/
func (this sliceValue) Slice(start, end int) (Value, bool) {
	if start < 0 {
		start = len(this) + start
	}

	if end < 0 {
		end = len(this) + end
	}

	if start <= end && start >= 0 && end <= len(this) {
		return NewValue(this[start:end]), true
	}

	return MISSING_VALUE, false
}

/*
If the start index is -ve, as per the N1QL specs, add it to
the length to get the actual index (from the end). If it is
valid(+ve) then return a slice from start till the length
of the slice and a bool value true. If the indices are not
valid return a missing value and false.
*/
func (this sliceValue) SliceTail(start int) (Value, bool) {
	if start < 0 {
		start = len(this) + start
	}

	if start >= 0 {
		return NewValue(this[start:]), true
	}

	return MISSING_VALUE, false
}

/*
It flattens out the elements of the array and appends it into
the buffer. This is done in child first (depth first) order.
In the event the buffer is full (capacity < length of the
buffer + the current element), then grow the buffer by
twice of length of the buffer + this element + 1.  Once the
buffer has space,range over the slice, append the children
to the buffer, and call Descendants recursively until there
are no elements left. Finally return the buffer.
*/
func (this sliceValue) Descendants(buffer []interface{}) []interface{} {
	if cap(buffer) < len(buffer)+len(this) {
		buf2 := make([]interface{}, len(buffer), (len(buffer)+len(this)+1)<<1)
		copy(buf2, buffer)
		buffer = buf2
	}

	for _, child := range this {
		buffer = append(buffer, child)
		buffer = NewValue(child).Descendants(buffer)
	}

	return buffer
}

/*
No fields to list. Hence return nil.
*/
func (this sliceValue) Fields() map[string]interface{} {
	return nil
}

/*
Append a small value.
*/
func (this sliceValue) Successor() Value {
	if len(this) == 0 {
		return _SMALL_ARRAY_VALUE
	}

	return sliceValue(append(this, nil))
}

func (this sliceValue) unwrap() Value {
	return this
}

var _SMALL_ARRAY_VALUE = sliceValue([]interface{}{nil})

/*
It is a struct containing slice values. This enables us to call all
the implemented methods for slicevalue without having to redefine them.
*/
type listValue struct {
	slice sliceValue
}

/*
Call implemented MarshalJSON method for slice in *listValue.
*/
func (this *listValue) MarshalJSON() ([]byte, error) {
	return this.slice.MarshalJSON()
}

/*
Type ARRAY.
*/
func (this *listValue) Type() Type { return ARRAY }

func (this *listValue) Actual() interface{} {
	return this.slice.Actual()
}

func (this *listValue) Equals(other Value) Value {
	return this.slice.Equals(other)
}

func (this *listValue) Collate(other Value) int {
	return this.slice.Collate(other)
}

func (this *listValue) Compare(other Value) Value {
	return this.slice.Compare(other)
}

func (this *listValue) Truth() bool {
	return this.slice.Truth()
}

/*
Call implemented Copy method for slice in *listValue.
Return a pointer to listValue whose entry is the return
value of the call to slicevalues copy method.
*/
func (this *listValue) Copy() Value {
	return &listValue{this.slice.Copy().(sliceValue)}
}

/*
Call implemented CopyForUpdate method for slice in *listValue.
*/
func (this *listValue) CopyForUpdate() Value {
	return this.slice.CopyForUpdate()
}

/*
Call implemented Field method for slice in *listValue.
*/
func (this *listValue) Field(field string) (Value, bool) {
	return this.slice.Field(field)
}

/*
Call implemented SetField method for slice in *listValue.
*/
func (this *listValue) SetField(field string, val interface{}) error {
	return this.slice.SetField(field, val)
}

/*
Call implemented UnsetField method for slice in *listValue.
*/
func (this *listValue) UnsetField(field string) error {
	return this.slice.UnsetField(field)
}

/*
Call implemented Index method for slice in *listValue.
*/
func (this *listValue) Index(index int) (Value, bool) {
	return this.slice.Index(index)
}

/*
It checks to see if there is a necessity to grow the slice.
If the index is greater than the length of the receiver
slice, check capacity next. In the event the index is
smaller than the capacity, assign the current slice to
the new slice from 0 to index+1. If the capacity is reached,
then grow the slice. Make a slice with length index+1 and
capacity twice the length, and reset the receiver. Finally
call the SetIndex method for the sliceValue.
*/
func (this *listValue) SetIndex(index int, val interface{}) error {
	if index >= len(this.slice) {
		if index < cap(this.slice) {
			this.slice = this.slice[0 : index+1]
		} else {
			slice := make(sliceValue, index+1, (index+1)<<1)
			copy(slice, this.slice)
			this.slice = slice
		}
	}

	return this.slice.SetIndex(index, val)
}

/*
Call implemented Slice method for slice in *listValue.
*/
func (this *listValue) Slice(start, end int) (Value, bool) {
	return this.slice.Slice(start, end)
}

/*
Call implemented SliceTail method for slice in *listValue.
*/
func (this *listValue) SliceTail(start int) (Value, bool) {
	return this.slice.SliceTail(start)
}

/*
Call implemented Descendants method for slice in *listValue.
*/
func (this *listValue) Descendants(buffer []interface{}) []interface{} {
	return this.slice.Descendants(buffer)
}

/*
Call implemented Fields method for slice in *listValue.
*/
func (this *listValue) Fields() map[string]interface{} {
	return this.slice.Fields()
}

/*
Append a small value.
*/
func (this *listValue) Successor() Value {
	return this.slice.Successor()
}

func (this *listValue) unwrap() Value {
	return this
}

func arrayEquals(array1, array2 []interface{}) Value {
	if len(array1) != len(array2) {
		return FALSE_VALUE
	}

	var missing, null Value
	for i, item1 := range array1 {
		eq := NewValue(item1).Equals(NewValue(array2[i]))
		switch eq.Type() {
		case MISSING:
			missing = eq
		case NULL:
			null = eq
		default:
			if !eq.Truth() {
				return eq
			}
		}
	}

	if missing != nil {
		return missing
	} else if null != nil {
		return null
	} else {
		return TRUE_VALUE
	}
}

/*
This code originally taken from https://github.com/couchbaselabs/walrus
*/
func arrayCollate(array1, array2 []interface{}) int {
	for i, item1 := range array1 {
		if i >= len(array2) {
			return 1
		}

		cmp := NewValue(item1).Collate(NewValue(array2[i]))
		if cmp != 0 {
			return cmp
		}
	}

	return len(array1) - len(array2)
}

func arrayCompare(array1, array2 []interface{}) Value {
	for i, item1 := range array1 {
		if i >= len(array2) {
			return ONE_VALUE
		}

		cmp := NewValue(item1).Compare(NewValue(array2[i]))
		if !cmp.Equals(ZERO_VALUE).Truth() {
			return cmp
		}
	}

	return NewValue(len(array1) - len(array2))
}

/*
It allows for a copy of every element of the array by using a copyFunc.
*/
func copySlice(source []interface{}, copier copyFunc) []interface{} {
	if source == nil {
		return nil
	}

	result := make([]interface{}, len(source))
	for i, v := range source {
		result[i] = copier(v)
	}

	return result
}

func marshalArray(slice []interface{}) (b []byte, err error) {
	if slice == nil {
		return _NULL_BYTES, nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	buf.WriteString("[")

	for i, e := range slice {
		if i > 0 {
			buf.WriteString(",")
		}

		v := NewValue(e)
		b, err = json.Marshal(v)
		if err != nil {
			return
		}

		buf.Write(b)
	}

	buf.WriteString("]")
	return buf.Bytes(), nil
}
