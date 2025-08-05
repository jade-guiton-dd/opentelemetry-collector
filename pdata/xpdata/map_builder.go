package xpdata

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type UnsafeMapBuilder struct {
	state  internal.State
	values []otlpcommon.KeyValue
}

func (mb *UnsafeMapBuilder) EnsureCapacity(capacity int) {
	oldValues := mb.values
	if capacity <= cap(oldValues) {
		return
	}
	mb.values = make([]otlpcommon.KeyValue, len(oldValues), capacity)
	copy(mb.values, oldValues)
}

func (mb *UnsafeMapBuilder) PutEmpty(k string) pcommon.Value {
	mb.values = append(mb.values, otlpcommon.KeyValue{Key: k})
	return pcommon.Value(internal.NewValue(&mb.values[len(mb.values)-1].Value, &mb.state))
}

func (mb *UnsafeMapBuilder) IntoMap() pcommon.Map {
	return pcommon.Map(internal.NewMap(&mb.values, &mb.state))
}
