//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/value"
)

///////////////////////////////////////////////////
//
// UnnestPosition
//
///////////////////////////////////////////////////

/*
UNNEST_POSITION(expr)
*/
type UnnestPosition struct {
	UnaryFunctionBase
}

func NewUnnestPosition(operand Expression) Function {
	rv := &UnnestPosition{
		*NewUnaryFunctionBase("unnest_position", operand),
	}

	rv.expr = rv
	return rv
}

func (this *UnnestPosition) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *UnnestPosition) Type() value.Type { return value.NUMBER }

func (this *UnnestPosition) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

func (this *UnnestPosition) Apply(context Context, operand value.Value) (value.Value, error) {
	if operand.Type() == value.MISSING {
		return operand, nil
	}

	av, ok := operand.(value.AnnotatedValue)
	if !ok {
		return value.NULL_VALUE, nil
	}

	upos := av.GetAttachment("unnest_position")
	if upos == nil {
		return value.NULL_VALUE, nil
	}

	pos, ok := upos.(int)
	if !ok {
		return value.NULL_VALUE, errors.NewUnnestInvalidPosition(pos)
	}

	return value.NewValue(pos), nil
}

func (this *UnnestPosition) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewUnnestPosition(operands[0])
	}
}
