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

	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/timestamp"
	"github.com/couchbase/query/util"
)

// Build a query execution pipeline from a query plan.
func Build(plan plan.Operator, context *Context) (Operator, error) {
	var m map[scannedIndex]bool
	if context.ScanVectorSource().Type() == timestamp.ONE_VECTOR {
		// Collect scanned indexes.
		m = make(map[scannedIndex]bool)
	}
	builder := &builder{context, m}
	x, err := plan.Accept(builder)

	if err != nil {
		return nil, err
	}

	if builder.scannedIndexes != nil && len(builder.scannedIndexes) > 1 {
		scannedIndexArr := make([]string, len(builder.scannedIndexes))
		for si := range builder.scannedIndexes {
			scannedIndexArr = append(scannedIndexArr, fmt.Sprintf("%s:%s", si.namespace, si.keyspace))
		}
		return nil, errors.NewScanVectorTooManyScannedBuckets(scannedIndexArr)
	}

	ex := x.(Operator)
	return ex, nil
}

type scannedIndex struct {
	namespace string
	keyspace  string
}

type builder struct {
	context        *Context
	scannedIndexes map[scannedIndex]bool // Nil if scanned indexes should not be collected.
}

// Scan
func (this *builder) VisitPrimaryScan(plan *plan.PrimaryScan) (interface{}, error) {
	panic("Unexpected VisitPrimaryScan")
	return nil, nil
}

func (this *builder) VisitParentScan(plan *plan.ParentScan) (interface{}, error) {
	panic("Unexpected VisitParentScan")
	return nil, nil
}

func (this *builder) VisitIndexScan(plan *plan.IndexScan) (interface{}, error) {
	return NewIndexScan(plan), nil
}

func (this *builder) VisitIndexCountScan(plan *plan.IndexCountScan) (interface{}, error) {
	panic("Unexpected VisitIndexCountScan")
	return nil, nil
}

func (this *builder) VisitKeyScan(plan *plan.KeyScan) (interface{}, error) {
	panic("Unexpected VisitKeyScan")
	return nil, nil
}

func (this *builder) VisitValueScan(plan *plan.ValueScan) (interface{}, error) {
	panic("Unexpected VisitValueScan")
	return nil, nil
}

func (this *builder) VisitDummyScan(plan *plan.DummyScan) (interface{}, error) {
	return NewDummyScan(), nil
}

func (this *builder) VisitCountScan(plan *plan.CountScan) (interface{}, error) {
	panic("Unexpected VisitCountScan")
	return nil, nil
}

func (this *builder) VisitIntersectScan(plan *plan.IntersectScan) (interface{}, error) {
	panic("Unexpected VisitIntersectScan")
	return nil, nil
}

func (this *builder) VisitUnionScan(plan *plan.UnionScan) (interface{}, error) {
	panic("Unexpected VisitUnionScan")
	return nil, nil
}

func (this *builder) VisitDistinctScan(plan *plan.DistinctScan) (interface{}, error) {
	panic("Unexpected VisitDistinctScan")
	return nil, nil
}

// Fetch
func (this *builder) VisitFetch(plan *plan.Fetch) (interface{}, error) {
	panic("Unexpected VisitFetch")
	return nil, nil
}

// DummyFetch
func (this *builder) VisitDummyFetch(plan *plan.DummyFetch) (interface{}, error) {
	panic("Unexpected VisitDummyFetch")
	return nil, nil
}

// Join
func (this *builder) VisitJoin(plan *plan.Join) (interface{}, error) {
	panic("Unexpected VisitJoin")
	return nil, nil
}

func (this *builder) VisitIndexJoin(plan *plan.IndexJoin) (interface{}, error) {
	panic("Unexpected VisitIndexJoin")
	return nil, nil
}

func (this *builder) VisitNest(plan *plan.Nest) (interface{}, error) {
	panic("Unexpected VisitNest")
	return nil, nil
}

func (this *builder) VisitIndexNest(plan *plan.IndexNest) (interface{}, error) {
	panic("Unexpected VisitIndexNest")
	return nil, nil
}

func (this *builder) VisitUnnest(plan *plan.Unnest) (interface{}, error) {
	panic("Unexpected VisitUnnest")
	return nil, nil
}

// Let + Letting
func (this *builder) VisitLet(plan *plan.Let) (interface{}, error) {
	panic("Unexpected VisitLet")
	return nil, nil
}

// Filter
func (this *builder) VisitFilter(plan *plan.Filter) (interface{}, error) {
	return NewFilter(plan), nil
}

// Group
func (this *builder) VisitInitialGroup(plan *plan.InitialGroup) (interface{}, error) {
	panic("Unexpected VisitInitialGroup")
	return nil, nil
}

func (this *builder) VisitIntermediateGroup(plan *plan.IntermediateGroup) (interface{}, error) {
	panic("Unexpected VisitIntermediateGroup")
	return nil, nil
}

func (this *builder) VisitFinalGroup(plan *plan.FinalGroup) (interface{}, error) {
	panic("Unexpected VisitFinalGroup")
	return nil, nil
}

// Project
func (this *builder) VisitInitialProject(plan *plan.InitialProject) (interface{}, error) {
	return NewInitialProject(plan), nil
}

func (this *builder) VisitFinalProject(plan *plan.FinalProject) (interface{}, error) {
	return NewFinalProject(), nil
}

func (this *builder) VisitIndexCountProject(plan *plan.IndexCountProject) (interface{}, error) {
	panic("Unexpected VisitIndexCountProject")
	return nil, nil
}

// Distinct
func (this *builder) VisitDistinct(plan *plan.Distinct) (interface{}, error) {
	panic("Unexpected VisitDistinct")
	return nil, nil
}

// Set operators
func (this *builder) VisitUnionAll(plan *plan.UnionAll) (interface{}, error) {
	panic("Unexpected VisitUnionAll")
	return nil, nil
}

func (this *builder) VisitIntersectAll(plan *plan.IntersectAll) (interface{}, error) {
	panic("Unexpected VisitIntersectAll")
	return nil, nil
}

func (this *builder) VisitExceptAll(plan *plan.ExceptAll) (interface{}, error) {
	panic("Unexpected VisitExceptAll")
	return nil, nil
}

// Order
func (this *builder) VisitOrder(plan *plan.Order) (interface{}, error) {
	panic("Unexpected VisitOrder")
	return nil, nil
}

// Offset
func (this *builder) VisitOffset(plan *plan.Offset) (interface{}, error) {
	panic("Unexpected VisitOffset")
	return nil, nil
}

func (this *builder) VisitLimit(plan *plan.Limit) (interface{}, error) {
	return NewLimit(plan), nil
}

// Insert
func (this *builder) VisitSendInsert(plan *plan.SendInsert) (interface{}, error) {
	panic("Unexpected VisitSendInsert")
	return nil, nil
}

// Upsert
func (this *builder) VisitSendUpsert(plan *plan.SendUpsert) (interface{}, error) {
	panic("Unexpected VisitSendUpsert")
	return nil, nil
}

// Delete
func (this *builder) VisitSendDelete(plan *plan.SendDelete) (interface{}, error) {
	panic("Unexpected VisitSendDelete")
	return nil, nil
}

// Update
func (this *builder) VisitClone(plan *plan.Clone) (interface{}, error) {
	panic("Unexpected VisitClone")
	return nil, nil
}

func (this *builder) VisitSet(plan *plan.Set) (interface{}, error) {
	panic("Unexpected VisitSet")
	return nil, nil
}

func (this *builder) VisitUnset(plan *plan.Unset) (interface{}, error) {
	panic("Unexpected VisitUnset")
	return nil, nil
}

func (this *builder) VisitSendUpdate(plan *plan.SendUpdate) (interface{}, error) {
	panic("Unexpected VisitSendUpdate")
	return nil, nil
}

// Merge
func (this *builder) VisitMerge(plan *plan.Merge) (interface{}, error) {
	panic("Unexpected VisitMerge")
	return nil, nil
}

// Alias
func (this *builder) VisitAlias(plan *plan.Alias) (interface{}, error) {
	panic("Unexpected VisitAlias")
	return nil, nil
}

// Authorize
func (this *builder) VisitAuthorize(plan *plan.Authorize) (interface{}, error) {
	child, err := plan.Child().Accept(this)
	if err != nil {
		return nil, err
	}

	return NewAuthorize(plan, child.(Operator)), nil
}

// Parallel
func (this *builder) VisitParallel(plan *plan.Parallel) (interface{}, error) {
	child, err := plan.Child().Accept(this)
	if err != nil {
		return nil, err
	}

	maxParallelism := util.MinInt(plan.MaxParallelism(), this.context.MaxParallelism())

	if maxParallelism == 1 {
		return child, nil
	} else {
		return NewParallel(plan, child.(Operator)), nil
	}
}

// Sequence
func (this *builder) VisitSequence(plan *plan.Sequence) (interface{}, error) {
	children := _SEQUENCE_POOL.Get()

	for _, pchild := range plan.Children() {
		child, err := pchild.Accept(this)
		if err != nil {
			return nil, err
		}

		children = append(children, child.(Operator))
	}

	return NewSequence(children...), nil
}

// Discard
func (this *builder) VisitDiscard(plan *plan.Discard) (interface{}, error) {
	panic("Unexpected VisitDiscard")
	return nil, nil
}

// Stream
func (this *builder) VisitStream(plan *plan.Stream) (interface{}, error) {
	return NewStream(), nil
}

// Collect
func (this *builder) VisitCollect(plan *plan.Collect) (interface{}, error) {
	panic("Unexpected VisitCollect")
	return nil, nil
}

// Channel
func (this *builder) VisitChannel(plan *plan.Channel) (interface{}, error) {
	panic("Unexpected VisitChannel")
	return nil, nil
}

// CreateIndex
func (this *builder) VisitCreatePrimaryIndex(plan *plan.CreatePrimaryIndex) (interface{}, error) {
	panic("Unexpected VisitCreatePrimaryIndex")
	return nil, nil
}

// CreateIndex
func (this *builder) VisitCreateIndex(plan *plan.CreateIndex) (interface{}, error) {
	panic("Unexpected VisitCreateIndex")
	return nil, nil
}

// DropIndex
func (this *builder) VisitDropIndex(plan *plan.DropIndex) (interface{}, error) {
	panic("Unexpected VisitDropIndex")
	return nil, nil
}

// AlterIndex
func (this *builder) VisitAlterIndex(plan *plan.AlterIndex) (interface{}, error) {
	panic("Unexpected VisitAlterIndex")
	return nil, nil
}

// BuildIndexes
func (this *builder) VisitBuildIndexes(plan *plan.BuildIndexes) (interface{}, error) {
	panic("Unexpected VisitBuildIndexes")
	return nil, nil
}

// Prepare
func (this *builder) VisitPrepare(plan *plan.Prepare) (interface{}, error) {
	return NewPrepare(plan.Prepared()), nil
}

// Explain
func (this *builder) VisitExplain(plan *plan.Explain) (interface{}, error) {
	panic("Unexpected VisitExplain")
	return nil, nil
}

// Infer
func (this *builder) VisitInferKeyspace(plan *plan.InferKeyspace) (interface{}, error) {
	panic("Unexpected VisitInferKeyspace")
	return nil, nil
}
