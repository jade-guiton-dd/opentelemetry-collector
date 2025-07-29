// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/collector/pdata/internal"

import (
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

type MapMap struct {
	orig  *map[string]*otlpcommon.AnyValue
	state *State
}

func GetOrigMapMap(ms MapMap) *map[string]*otlpcommon.AnyValue {
	return ms.orig
}

func GetMapMapState(ms MapMap) *State {
	return ms.state
}

func NewMapMap(orig *map[string]*otlpcommon.AnyValue, state *State) MapMap {
	return MapMap{orig: orig, state: state}
}
