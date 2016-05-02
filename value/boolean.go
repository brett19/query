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
	"math"

	"github.com/couchbase/query/util"
	"bytes"
)

/*
BoolValue is defined as a bool type.
*/
type boolValue bool

/*
FALSE_VALUE/ TRUE_VALUE are assigned a false and true value
(NewValue() is called to ensure that it is a value), and
_FALSE _BYTES / _TRUE _BYTES that are slices of bytes
representing false and true.
*/
var FALSE_VALUE = NewValue(false)
var TRUE_VALUE = NewValue(true)

/*
_FALSE _BYTES / _TRUE _BYTES that are slices of bytes
representing false and true.
*/
var _FALSE_BYTES = []byte("false")
var _TRUE_BYTES = []byte("true")

func (this boolValue) String() string {
	if this {
		return "true"
	} else {
		return "false"
	}
}

func (this boolValue) MarshalJSON() ([]byte, error) {
	if this {
		return _TRUE_BYTES, nil
	} else {
		return _FALSE_BYTES, nil
	}
}

func (this boolValue) FastMarshalJSON(buf *bytes.Buffer) error {
	if this {
		buf.Write(_TRUE_BYTES)
	} else {
		buf.Write(_FALSE_BYTES)
	}
	return nil
}

/*
Type BOOLEAN
*/
func (this boolValue) Type() Type {
	return BOOLEAN
}

/*
Cast receiver to bool and return.
*/
func (this boolValue) Actual() interface{} {
	return bool(this)
}

/*
If other is a boolValue, compare it with the receiver. If
it is a parsedValue or annotated value then call Equals
by parsing other or Values respectively. If it is any other
type we return false.
*/
func (this boolValue) Equals(other Value) Value {
	other = other.unwrap()
	switch other := other.(type) {
	case missingValue:
		return other
	case *nullValue:
		return other
	case boolValue:
		if this == other {
			return TRUE_VALUE
		}
	}

	return FALSE_VALUE
}

func (this boolValue) Collate(other Value) int {
	other = other.unwrap()
	switch other := other.(type) {
	case boolValue:
		if this == other {
			return 0
		} else if !this {
			return -1
		} else {
			return 1
		}
	default:
		return int(BOOLEAN - other.Type())
	}
}

func (this boolValue) Compare(other Value) Value {
	other = other.unwrap()
	switch other := other.(type) {
	case missingValue:
		return other
	case *nullValue:
		return other
	default:
		return NewValue(this.Collate(other))
	}
}

/*
Cast receiver to bool and return.
*/
func (this boolValue) Truth() bool {
	return bool(this)
}

/*
Return receiver.
*/
func (this boolValue) Copy() Value {
	return this
}

/*
Return receiver.
*/
func (this boolValue) CopyForUpdate() Value {
	return this
}

/*
Calls missingField.
*/
func (this boolValue) Field(field string) (Value, bool) {
	return missingField(field), false
}

/*
Not valid for bool.
*/
func (this boolValue) SetField(field string, val interface{}) error {
	return Unsettable(field)
}

/*
Not valid for bool.
*/
func (this boolValue) UnsetField(field string) error {
	return Unsettable(field)
}

/*
Calls missingIndex.
*/
func (this boolValue) Index(index int) (Value, bool) {
	return missingIndex(index), false
}

/*
Not valid for bool.
*/
func (this boolValue) SetIndex(index int, val interface{}) error {
	return Unsettable(index)
}

/*
Returns NULL_VALUE
*/
func (this boolValue) Slice(start, end int) (Value, bool) {
	return NULL_VALUE, false
}

/*
Returns NULL_VALUE
*/
func (this boolValue) SliceTail(start int) (Value, bool) {
	return NULL_VALUE, false
}

/*
Returns the input buffer as is.
*/
func (this boolValue) Descendants(buffer []interface{}) []interface{} {
	return buffer
}

/*
Bool has no fields to list. Hence return nil.
*/
func (this boolValue) Fields() map[string]interface{} {
	return nil
}

func (this boolValue) FieldNames(buffer []string) []string {
	return nil
}

/*
Returns the input buffer as is.
*/
func (this boolValue) DescendantPairs(buffer []util.IPair) []util.IPair {
	return buffer
}

/*
FALSE is succeeded by TRUE, TRUE by numbers.
*/
func (this boolValue) Successor() Value {
	if bool(this) {
		return _MIN_NUMBER_VALUE
	} else {
		return TRUE_VALUE
	}
}

func (this boolValue) unwrap() Value {
	return this
}

var _MIN_NUMBER_VALUE = NewValue(-math.MaxFloat64)
