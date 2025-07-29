// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/collector/pdata/internal"

import (
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

type SortedMap struct {
	orig  *[]otlpcommon.KeyValue
	state *State
}

func GetOrigSortedMap(ms SortedMap) *[]otlpcommon.KeyValue {
	return ms.orig
}

func GetSortedMapState(ms SortedMap) *State {
	return ms.state
}

func NewSortedMap(orig *[]otlpcommon.KeyValue, state *State) SortedMap {
	return SortedMap{orig: orig, state: state}
}
