// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pcommon // import "go.opentelemetry.io/collector/pdata/pcommon"

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

// MapMap stores a map of string keys to elements of Value type.
//
// Must use NewMapMap function to create new instances.
// Important: zero-initialized instance is not valid for use.
type MapMap internal.MapMap

// NewMapMap creates a MapMap with 0 elements.
func NewMapMap() MapMap {
	orig := map[string]*otlpcommon.AnyValue{}
	state := internal.StateMutable
	return MapMap(internal.NewMapMap(&orig, &state))
}

func (m MapMap) getOrig() *map[string]*otlpcommon.AnyValue {
	return internal.GetOrigMapMap(internal.MapMap(m))
}

func (m MapMap) getState() *internal.State {
	return internal.GetMapMapState(internal.MapMap(m))
}

// EnsureCapacity increases the capacity of this MapMap instance, if necessary,
// to ensure that it can hold at least the number of elements specified by the capacity argument.
func (m MapMap) EnsureCapacity(capacity int) {
	m.getState().AssertMutable()

	oldOrig := *m.getOrig()
	// so you can create a map with a capacity, but not check the existing capacity??
	newOrig := make(map[string]*otlpcommon.AnyValue, capacity)
	for k, v := range oldOrig {
		newOrig[k] = v
	}
	*m.getOrig() = newOrig
}

// PutEmpty inserts or updates an empty value to the map under given key
// and return the updated/inserted value.
func (m MapMap) PutEmpty(k string) Value {
	m.getState().AssertMutable()
	if val, existing := (*m.getOrig())[k]; existing {
		av := newValue(val, m.getState())
		av.getOrig().Value = nil
		return newValue(av.getOrig(), m.getState())
	}
	val := &otlpcommon.AnyValue{}
	(*m.getOrig())[k] = val
	return newValue(val, m.getState())
}
