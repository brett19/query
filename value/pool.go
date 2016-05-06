package value

import (
	"sync"
)

type Pool struct {
	flatObjectValuePool sync.Pool
	objectValuePool sync.Pool
	annotatedValuePool sync.Pool
	scopeValuePool sync.Pool
	annotatedScopeValuePool sync.Pool
}

func (p *Pool) NewFlatObjectValue(initialSize int) *flatObject {
	pv := p.flatObjectValuePool.Get()
	var v *flatObject
	if pv == nil {
		v = &flatObject{}
	} else {
		v = pv.(*flatObject)
	}
	return v
}

func (p *Pool) ReleaseFlatObjectValue(val *flatObject) {
	val.items.Clear()
	val.parsed = nil
	p.flatObjectValuePool.Put(val)
}

func (p *Pool) NewEmptyObjectValue(initialSize int) *objectValue {
	pv := p.objectValuePool.Get()
	var v *objectValue
	if pv == nil {
		nv := objectValue(make(map[string]interface{}, initialSize))
		v = &nv
	} else {
		v = pv.(*objectValue)
	}
	return v
}

func (p *Pool) ReleaseObjectValue(val *objectValue) {
	for i, _ := range *val {
		delete(*val, i)
	}
	p.objectValuePool.Put(val)
}

func (p *Pool) NewAnnotatedValue(val interface{}) *annotatedValue {
	switch val := val.(type) {
	case *annotatedValue:
		return val
	case Value:
		return p.newAnnotatedValue(val)
	default:
		return p.NewAnnotatedValue(NewValue(val))
	}
}

func (p *Pool) newAnnotatedValue(val Value) *annotatedValue {
	pv := p.annotatedValuePool.Get()
	var v *annotatedValue
	if pv == nil {
		v = &annotatedValue{}
	} else {
		v = pv.(*annotatedValue)
	}
	v.Value = val
	return v
}

func (p *Pool) ReleaseAnnotatedValue(val *annotatedValue) {
	p.Release(val.Value)
	val.Value = nil
	for i, v := range val.attachments {
		p.Release(v)
		delete(val.attachments, i)
	}
	for i, _ := range val.covers {
		delete(val.covers, i)
	}
	p.annotatedValuePool.Put(val)
}

func (p *Pool) NewScopeValue(val Value, parent Value) *ScopeValue {
	pv := p.scopeValuePool.Get()
	var v *ScopeValue
	if pv == nil {
		v = &ScopeValue{}
	} else {
		v = pv.(*ScopeValue)
	}

	v.Value = val
	v.parent = parent
	return v
}

func (p *Pool) ReleaseScopeValue(val *ScopeValue) {
	p.Release(val.Value)
	p.Release(val.parent)
	val.Value = nil
	val.parent = nil
	p.scopeValuePool.Put(val)
}

func (p *Pool) NewAnnotatedScopeValue(val Value, parent Value) *annotatedScopeValue {
	pv := p.annotatedScopeValuePool.Get()
	var v *annotatedScopeValue
	if pv == nil {
		v = &annotatedScopeValue{}
	} else {
		v = pv.(*annotatedScopeValue)
	}

	v.Value = val
	v.parent = p.NewAnnotatedValue(parent)
	return v
}

func (p *Pool) ReleaseAnnotatedScopeValue(val *annotatedScopeValue) {
	p.Release(val.Value)
	p.Release(val.parent)
	val.Value = nil
	val.parent = nil
	p.annotatedScopeValuePool.Put(val)
}

func (p *Pool) Release(item interface{}) {
	if item == nil {
		return
	}
	switch v := item.(type) {
	case *flatObject:
		p.ReleaseFlatObjectValue(v)
	case *objectValue:
		p.ReleaseObjectValue(v)
	case *annotatedValue:
		p.ReleaseAnnotatedValue(v)
	case *annotatedScopeValue:
		p.ReleaseAnnotatedScopeValue(v)
	case *ScopeValue:
		p.ReleaseScopeValue(v)
	}
}
