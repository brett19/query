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
	"sync"
	go_atomic "sync/atomic"

	atomic "github.com/couchbase/go-couchbase/platform"
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/value"
)

type base struct {
	input       Operator
	output      Operator
	once        sync.Once
	batch       []value.AnnotatedValue
}

const _ITEM_CAP = 512

var pipelineCap atomic.AlignedInt64

func init() {
	atomic.StoreInt64(&pipelineCap, int64(_ITEM_CAP))
	p := value.NewAnnotatedPool(_BATCH_SIZE)
	_BATCH_POOL.Store(p)
}

func SetPipelineCap(cap int) {
	if cap < 1 {
		cap = _ITEM_CAP
	}
	atomic.StoreInt64(&pipelineCap, int64(cap))
}

func GetPipelineCap() int64 {
	return atomic.LoadInt64(&pipelineCap)
}

func newBase() base {
	return base{
	}
}

// The output of this operator will be redirected elsewhere, so we
// allocate a minimal itemChannel.
func newRedirectBase() base {
	return base{
	}
}

func (this *base) Input() Operator {
	return this.input
}

func (this *base) SetInput(op Operator) {
	this.input = op
}

func (this *base) Output() Operator {
	return this.output
}

func (this *base) SetOutput(op Operator) {
	this.output = op
}

func (this *base) Item(item value.AnnotatedValue, context *Context) bool {
	panic("item lost due to unimplemented handler")
	return true
}

func (this *base) copy() base {
	return base{
		input:       this.input,
		output:      this.output,
	}
}

func (this *base) sendItem(item value.AnnotatedValue, context *Context) bool {
	this.output.Item(item, context)
	return true
}

type consumer interface {
	Readonly() bool
	BeforeItems(context *Context, parent value.Value) bool
	AfterItems(context *Context)
}

func (this *base) Readonly() bool {
	return true
}

func (this *base) BeforeItems(context *Context, parent value.Value) bool {
	return true
}

func (this *base) AfterItems(context *Context) {
}

func (this *base) runConsumer(cons consumer, context *Context, parent value.Value) {
	if context.Readonly() && !cons.Readonly() {
		return
	}

	ok := cons.BeforeItems(context, parent)

	if ok {
		this.input.RunOnce(context, parent)
	}

	cons.AfterItems(context)
}

type batcher interface {
	allocateBatch()
	enbatch(item value.AnnotatedValue, b batcher, context *Context) bool
	enbatchSize(item value.AnnotatedValue, b batcher, batchSize int, context *Context) bool
	flushBatch(context *Context) bool
	releaseBatch()
}

var _BATCH_SIZE = 64

var _BATCH_POOL go_atomic.Value

func SetPipelineBatch(size int) {
	if size < 1 {
		size = _BATCH_SIZE
	}

	p := value.NewAnnotatedPool(size)
	_BATCH_POOL.Store(p)
}

func PipelineBatchSize() int {
	return _BATCH_POOL.Load().(*value.AnnotatedPool).Size()
}

func getBatchPool() *value.AnnotatedPool {
	return _BATCH_POOL.Load().(*value.AnnotatedPool)
}

func (this *base) allocateBatch() {
	this.batch = getBatchPool().Get()
}

func (this *base) releaseBatch() {
	getBatchPool().Put(this.batch)
	this.batch = nil
}

func (this *base) enbatchSize(item value.AnnotatedValue, b batcher, batchSize int, context *Context) bool {
	if this.batch == nil {
		this.allocateBatch()
	} else if len(this.batch) == batchSize {
		if !b.flushBatch(context) {
			return false
		}

		if len(this.batch) == batchSize {
			this.allocateBatch()
		}
	}

	this.batch = append(this.batch, item)
	return true
}
func (this *base) enbatch(item value.AnnotatedValue, b batcher, context *Context) bool {
	return this.enbatchSize(item, b, cap(this.batch), context)
}

func (this *base) requireKey(item value.AnnotatedValue, context *Context) (string, bool) {
	mv := item.GetAttachment("meta")
	if mv == nil {
		context.Error(errors.NewInvalidValueError(
			fmt.Sprintf("Value does not contain META: %v", item)))
		return "", false
	}

	meta := mv.(map[string]interface{})
	key, ok := meta["id"]
	if !ok {
		context.Error(errors.NewInvalidValueError(
			fmt.Sprintf("META does not contain ID: %v", item)))
		return "", false
	}

	act := value.NewValue(key).Actual()
	switch act := act.(type) {
	case string:
		return act, true
	default:
		context.Error(errors.NewInvalidValueError(
			fmt.Sprintf("ID %v of type %T is not a string in value %v", act, act, item)))
		return "", false
	}
}
