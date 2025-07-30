package pcommon

import (
	"cmp"
	"slices"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

func (m Map) PutEmptyUnsafe(k string) Value {
	*m.getOrig() = append(*m.getOrig(), otlpcommon.KeyValue{Key: k})
	return newValue(&(*m.getOrig())[len(*m.getOrig())-1].Value, m.getState())
}

func (m Map) PutEmptyLessUnsafe(k string) Value {
	m.getState().AssertMutable()
	*m.getOrig() = append(*m.getOrig(), otlpcommon.KeyValue{Key: k})
	return newValue(&(*m.getOrig())[len(*m.getOrig())-1].Value, m.getState())
}

type SortedMapBuilder struct {
	state  internal.State
	values []otlpcommon.KeyValue
}

func (mb *SortedMapBuilder) EnsureCapacity(capacity int) {
	oldValues := mb.values
	if capacity <= cap(oldValues) {
		return
	}
	mb.values = make([]otlpcommon.KeyValue, len(oldValues), capacity)
	copy(mb.values, oldValues)
}

func (mb *SortedMapBuilder) PutEmpty(k string) Value {
	if len(mb.values) > 0 && mb.values[len(mb.values)-1].Key >= k {
		panic("keys added to SortedMapBuilder were not in strictly ascending order")
	}
	mb.values = append(mb.values, otlpcommon.KeyValue{Key: k})
	return newValue(&mb.values[len(mb.values)-1].Value, &mb.state)
}

func (mb *SortedMapBuilder) IntoMap() Map {
	return newMap(&mb.values, &mb.state)
}

type MapBuilder struct {
	state  internal.State
	values []otlpcommon.KeyValue
}

func (mb *MapBuilder) EnsureCapacity(capacity int) {
	oldValues := mb.values
	if capacity <= cap(oldValues) {
		return
	}
	mb.values = make([]otlpcommon.KeyValue, len(oldValues), capacity)
	copy(mb.values, oldValues)
}

func (mb *MapBuilder) PutEmpty(k string) Value {
	mb.values = append(mb.values, otlpcommon.KeyValue{Key: k})
	return newValue(&mb.values[len(mb.values)-1].Value, &mb.state)
}

func (mb *MapBuilder) IntoMap(merge func(dst Value, src Value, firstMerge bool)) Map {
	if len(mb.values) == 0 {
		return newMap(&mb.values, &mb.state)
	}
	slices.SortFunc(mb.values, func(kv1 otlpcommon.KeyValue, kv2 otlpcommon.KeyValue) int {
		return cmp.Compare(kv1.Key, kv2.Key)
	})
	from := 1
	last := 0
	firstMerge := true
	n := len(mb.values)
	for from < n {
		if mb.values[from].Key == mb.values[last].Key {
			merge(
				newValue(&mb.values[last].Value, &mb.state),
				newValue(&mb.values[from].Value, &mb.state),
				firstMerge,
			)
			firstMerge = false
		} else {
			last++
			firstMerge = true
			if from != last {
				mb.values[last] = mb.values[from]
			}
		}
		from++
	}
	mb.values = mb.values[:last+1]
	return newMap(&mb.values, &mb.state)
}
