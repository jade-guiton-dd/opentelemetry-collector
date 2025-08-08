package xpdata

import (
	"cmp"
	"slices"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type MapBuilder struct {
	state internal.State
	pairs []otlpcommon.KeyValue
}

func (mb *MapBuilder) EnsureCapacity(capacity int) {
	mb.state.AssertMutable()
	oldValues := mb.pairs
	if capacity <= cap(oldValues) {
		return
	}
	mb.pairs = make([]otlpcommon.KeyValue, len(oldValues), capacity)
	copy(mb.pairs, oldValues)
}

func compareKeyValues(kv1 otlpcommon.KeyValue, kv2 otlpcommon.KeyValue) int {
	return cmp.Compare(kv1.Key, kv2.Key)
}

func (mb *MapBuilder) getValue(i int) pcommon.Value {
	return pcommon.Value(internal.NewValue(&mb.pairs[i].Value, &mb.state))
}

func (mb *MapBuilder) AppendEmpty(k string) pcommon.Value {
	mb.state.AssertMutable()
	mb.pairs = append(mb.pairs, otlpcommon.KeyValue{Key: k})
	return mb.getValue(len(mb.pairs) - 1)
}

func (mb *MapBuilder) UnsafeIntoMap(m pcommon.Map) {
	mb.state.AssertMutable()
	internal.GetMapState(internal.Map(m)).AssertMutable()
	mb.state = internal.StateReadOnly // to avoid modifying a Map later marked as ReadOnly through builder Values
	*internal.GetOrigMap(internal.Map(m)) = mb.pairs
}

func (mb *MapBuilder) SortedIntoMap(m pcommon.Map) {
	for i := range len(mb.pairs) - 1 {
		if mb.pairs[i].Key >= mb.pairs[i+1].Key {
			panic("Keys added to MapBuilder are not in strictly ascending order")
		}
	}
	mb.UnsafeIntoMap(m)
}

func (mb *MapBuilder) DistinctIntoMap(m pcommon.Map) {
	slices.SortFunc(mb.pairs, compareKeyValues)
	for i := range len(mb.pairs) - 1 {
		if mb.pairs[i].Key == mb.pairs[i+1].Key {
			panic("Keys added to MapBuilder are not distinct")
		}
	}
	mb.UnsafeIntoMap(m)
}

func (mb *MapBuilder) MergeIntoMap(m pcommon.Map, merge func(vals []pcommon.Value)) {
	if len(mb.pairs) > 0 {
		slices.SortFunc(mb.pairs, compareKeyValues)
		groupStart := 0
		groupEnd := 1
		groupKey := mb.pairs[0].Key
		for {
			if groupEnd < len(mb.pairs) && mb.pairs[groupEnd].Key == groupKey {
				groupEnd++
				continue
			}
			groupSize := groupEnd - groupStart
			if groupSize > 1 {
				vals := make([]pcommon.Value, groupSize)
				for i := range groupSize {
					vals[i] = mb.getValue(groupStart + i)
				}
				// We expect the result to be placed in vals[0]
				merge(vals)
				mb.pairs = slices.Delete(mb.pairs, groupStart+1, groupEnd)
				groupEnd = groupStart + 1
			}
			if groupEnd == len(mb.pairs) {
				break
			}
			groupStart = groupEnd
			groupEnd = groupEnd + 1
			groupKey = mb.pairs[groupStart].Key
		}
	}
	mb.UnsafeIntoMap(m)
}
