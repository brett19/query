//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package planner

import (
	"github.com/couchbaselabs/query/datastore"
	"github.com/couchbaselabs/query/expression"
	"github.com/couchbaselabs/query/value"
)

type sargLT struct {
	sargBase
}

func newSargLT(expr *expression.LT) *sargLT {
	rv := &sargLT{}
	rv.sarg = func(expr2 expression.Expression) (datastore.Spans, error) {
		if expr.EquivalentTo(expr2) {
			return _SELF_SPANS, nil
		}

		var values value.Values
		span := &datastore.Span{}

		if expr.First().EquivalentTo(expr2) {
			values = value.Values{expr.Second().Value()}
			span.Range.High = values
		} else if expr.Second().EquivalentTo(expr2) {
			values = value.Values{expr.First().Value()}
			span.Range.Low = values
		}

		if len(values) == 0 {
			return nil, nil
		}

		span.Range.Inclusion = datastore.NEITHER
		return datastore.Spans{span}, nil
	}

	return rv
}