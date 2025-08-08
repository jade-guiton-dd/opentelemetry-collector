// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pcommon // import "go.opentelemetry.io/collector/pdata/pcommon"

import (
	"cmp"
	"iter"
	"slices"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

// SortedMap stores a map of string keys to elements of Value type.
//
// Must use NewSortedMap function to create new instances.
// Important: zero-initialized instance is not valid for use.
type SortedMap internal.SortedMap

// NewSortedMap creates a SortedMap with 0 elements.
func NewSortedMap() SortedMap {
	orig := []otlpcommon.KeyValue(nil)
	state := internal.StateMutable
	return SortedMap(internal.NewSortedMap(&orig, &state))
}

func (m SortedMap) getOrig() *[]otlpcommon.KeyValue {
	return internal.GetOrigSortedMap(internal.SortedMap(m))
}

func (m SortedMap) getState() *internal.State {
	return internal.GetSortedMapState(internal.SortedMap(m))
}

// EnsureCapacity increases the capacity of this SortedMap instance, if necessary,
// to ensure that it can hold at least the number of elements specified by the capacity argument.
func (m SortedMap) EnsureCapacity(capacity int) {
	m.getState().AssertMutable()
	oldOrig := *m.getOrig()
	if capacity <= cap(oldOrig) {
		return
	}
	*m.getOrig() = make([]otlpcommon.KeyValue, len(oldOrig), capacity)
	copy(*m.getOrig(), oldOrig)
}

func (m SortedMap) find(key string) (int, bool) {
	slice := *m.getOrig()

	if len(slice) < 8 {
		for i := range slice {
			key2 := slice[i].Key
			if key2 == key {
				return i, true
			} else if key2 > key {
				return i, false
			}
		}
		return len(slice), false
	}

	return slices.BinarySearchFunc(slice, key, func(kv otlpcommon.KeyValue, target string) int {
		return cmp.Compare(kv.Key, target)
	})
}

func (m SortedMap) Get(key string) (Value, bool) {
	if i, found := m.find(key); found {
		return newValue(&(*m.getOrig())[i].Value, m.getState()), true
	}
	return newValue(nil, m.getState()), false
}

// PutEmpty inserts or updates an empty value to the map under given key
// and return the updated/inserted value.
func (m SortedMap) PutEmpty(k string) Value {
	m.getState().AssertMutable()
	i, existing := m.find(k)
	if existing {
		av := newValue(&(*m.getOrig())[i].Value, m.getState())
		av.getOrig().Value = nil
		return newValue(av.getOrig(), m.getState())
	}
	*m.getOrig() = slices.Insert(*m.getOrig(), i, otlpcommon.KeyValue{Key: k})
	return newValue(&(*m.getOrig())[i].Value, m.getState())
}

func (m SortedMap) PutStr(k string, v string) {
	m.PutEmpty(k).SetStr(v)
}

func (m SortedMap) All() iter.Seq2[string, Value] {
	return func(yield func(string, Value) bool) {
		for i := range *m.getOrig() {
			kv := &(*m.getOrig())[i]
			if !yield(kv.Key, Value(internal.NewValue(&kv.Value, m.getState()))) {
				return
			}
		}
	}
}
