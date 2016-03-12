//  Copyright (c) 2016 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"github.com/couchbase/query/value"
)

/*
Represents range predicate Every, that allow testing of a bool
condition over the elements or attributes of a collection or
object. Type AnyEvery is a struct that implements collPred.
*/
type AnyEvery struct {
	collPredBase
}

/*
This method returns a pointer to the AnyEvery struct that has the
bindings and satisfies fields populated by the input args
bindings and expression satisfies.
*/
func NewAnyEvery(bindings Bindings, satisfies Expression) Expression {
	rv := &AnyEvery{
		collPredBase: collPredBase{
			bindings:  bindings,
			satisfies: satisfies,
		},
	}

	rv.expr = rv
	return rv
}

/*
It calls the VisitAnyEvery method by passing in the receiver to
and returns the interface. It is a visitor pattern.
*/
func (this *AnyEvery) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAnyEvery(this)
}

/*
It returns a Boolean value.
*/
func (this *AnyEvery) Type() value.Type { return value.BOOLEAN }

func (this *AnyEvery) Evaluate(item value.Value, context Context) (value.Value, error) {
	missing := false
	null := false

	barr := make([][]interface{}, len(this.bindings))
	for i, b := range this.bindings {
		bv, err := b.Expression().Evaluate(item, context)
		if err != nil {
			return nil, err
		}

		if b.Descend() {
			buffer := make([]interface{}, 0, 256)
			bv = value.NewValue(bv.Descendants(buffer))
		}

		switch bv.Type() {
		case value.ARRAY:
			barr[i] = bv.Actual().([]interface{})
		case value.MISSING:
			missing = true
		default:
			null = true
		}
	}

	if missing {
		return value.MISSING_VALUE, nil
	}

	if null {
		return value.NULL_VALUE, nil
	}

	n := -1
	for _, b := range barr {
		if n < 0 || len(b) < n {
			n = len(b)
		}
	}

	for i := 0; i < n; i++ {
		cv := value.NewScopeValue(make(map[string]interface{}, len(this.bindings)), item)
		for j, b := range this.bindings {
			cv.SetField(b.Variable(), barr[j][i])
		}

		sv, e := this.satisfies.Evaluate(cv, context)
		if e != nil {
			return nil, e
		}

		if !sv.Truth() {
			return value.NewValue(false), nil
		}
	}

	return value.NewValue(n > 0), nil
}

func (this *AnyEvery) Copy() Expression {
	return NewAnyEvery(this.bindings.Copy(), Copy(this.satisfies))
}