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
	"fmt"
	"math"

	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

type IndexScan struct {
	base
	plan         *plan.IndexScan
	outChannel   chan value.AnnotatedValue
}

func NewIndexScan(plan *plan.IndexScan) *IndexScan {
	rv := &IndexScan{
		base: newBase(),
		plan: plan,
		outChannel:   make(chan value.AnnotatedValue, 40),
	}

	rv.output = rv
	return rv
}

func (this *IndexScan) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitIndexScan(this)
}

func (this *IndexScan) Copy() Operator {
	return &IndexScan{
		base: this.base.copy(),
		plan: this.plan,
	}
}

func (this *IndexScan) Item(item value.AnnotatedValue, context *Context) bool {
	this.outChannel <- item
	return true
}

func (this *IndexScan) RunOnce(context *Context, parent value.Value) {
	context.AddPhaseOperator(INDEX_SCAN)
	defer context.Recover()

	spans := this.plan.Spans()
	n := len(spans)

	children := _INDEX_SCAN_POOL.Get()
	defer _INDEX_SCAN_POOL.Put(children)

	doneCh := make(chan bool)

	for _, span := range spans {
		span := span
		go func() {
			execSpan := newSpanScan(this, span)
			execSpan.SetOutput(this)
			execSpan.RunOnce(context, parent)
			doneCh <- true
		}()
	}

	go func() {
		for i := 0; i < n; i++ {
			<- doneCh
		}
		close(this.outChannel)
	}()

	Loop:
	for {
		item, ok := <-this.outChannel
		if !ok {
			break Loop
		}
		this.sendItem(item, context)
	}
}

type spanScan struct {
	base
	parent *IndexScan
	plan *plan.IndexScan
	span *plan.Span
}

func newSpanScan(parent *IndexScan, span *plan.Span) *spanScan {
	rv := &spanScan{
		base: newBase(),
		parent: parent,
		plan: parent.plan,
		span: span,
	}

	rv.output = parent.output
	return rv
}

func (this *spanScan) Accept(visitor Visitor) (interface{}, error) {
	panic(fmt.Sprintf("Internal operator spanScan visited by %v.", visitor))
}

func (this *spanScan) Copy() Operator {
	return &spanScan{
		base: this.base.copy(),
		parent: this.parent,
		plan: this.plan,
		span: this.span,
	}
}

func (this *spanScan) RunOnce(context *Context, parent value.Value) {
	defer context.Recover()

	conn := datastore.NewIndexConnection(context)
	defer notifyConn(conn.StopChannel())

	go this.scan(context, conn)

	var entry *datastore.IndexEntry
	ok := true
	var docs uint64 = 0

	var countDocs = func() {
		if docs > 0 {
			context.AddPhaseCount(INDEX_SCAN, docs)
		}
	}
	defer countDocs()

	for ok {
		select {
		case entry, ok = <-conn.EntryChannel():
			if ok {
				cv := value.NewScopeValue(make(map[string]interface{}), parent)
				av := value.NewAnnotatedValue(cv)

				// For downstream Fetch
				meta := map[string]interface{}{"id": entry.PrimaryKey}
				av.SetAttachment("meta", meta)

				covers := this.plan.Covers()
				if len(covers) > 0 {
					for c, v := range this.plan.FilterCovers() {
						av.SetCover(c.Text(), v)
					}

					// Matches planner.builder.buildCoveringScan()
					for i, ek := range entry.EntryKey {
						av.SetCover(covers[i].Text(), ek)
					}

					// Matches planner.builder.buildCoveringScan()
					av.SetCover(covers[len(entry.EntryKey)].Text(),
						value.NewValue(entry.PrimaryKey))

					av.SetField(this.plan.Term().Alias(), av)
				}

				ok = this.sendItem(av, context)
				docs++
				if docs > _PHASE_UPDATE_COUNT {
					context.AddPhaseCount(INDEX_SCAN, docs)
					docs = 0
				}
			}
		}
	}
}

func (this *spanScan) scan(context *Context, conn *datastore.IndexConnection) {
	defer context.Recover() // Recover from any panic

	dspan, err := evalSpan(this.span, context)
	if err != nil {
		context.Error(errors.NewEvaluationError(err, "span"))
		close(conn.EntryChannel())
		return
	}

	limit := int64(math.MaxInt64)
	if this.plan.Limit() != nil {
		if context.ScanConsistency() == datastore.UNBOUNDED || this.plan.Covers() != nil {
			lv, err := this.plan.Limit().Evaluate(nil, context)
			if err == nil && lv.Type() == value.NUMBER {
				limit = int64(lv.Actual().(float64))
			}
		}
	}

	keyspaceTerm := this.plan.Term()
	scanVector := context.ScanVectorSource().ScanVector(keyspaceTerm.Namespace(), keyspaceTerm.Keyspace())
	this.plan.Index().Scan(context.RequestId(), dspan, this.plan.Distinct(), limit,
		context.ScanConsistency(), scanVector, conn)
}

func evalSpan(ps *plan.Span, context *Context) (*datastore.Span, error) {
	var err error
	ds := &datastore.Span{}

	ds.Seek, err = evalExprs(ps.Seek, context)
	if err != nil {
		return nil, err
	}

	ds.Range.Low, err = evalExprs(ps.Range.Low, context)
	if err != nil {
		return nil, err
	}

	ds.Range.High, err = evalExprs(ps.Range.High, context)
	if err != nil {
		return nil, err
	}

	ds.Range.Inclusion = ps.Range.Inclusion
	return ds, nil
}

func evalExprs(exprs expression.Expressions, context *Context) (value.Values, error) {
	if exprs == nil {
		return nil, nil
	}

	values := make(value.Values, len(exprs))

	var err error
	for i, expr := range exprs {
		if expr == nil {
			continue
		}

		values[i], err = expr.Evaluate(nil, context)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}