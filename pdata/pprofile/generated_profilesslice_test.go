// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpprofiles "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1development"
)

func TestProfilesSlice(t *testing.T) {
	es := NewProfilesSlice()
	assert.Equal(t, 0, es.Len())
	state := internal.StateMutable
	es = newProfilesSlice(&[]*otlpprofiles.Profile{}, &state)
	assert.Equal(t, 0, es.Len())

	emptyVal := NewProfile()
	testVal := generateTestProfile()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestProfile(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestProfilesSliceReadOnly(t *testing.T) {
	sharedState := internal.StateReadOnly
	es := newProfilesSlice(&[]*otlpprofiles.Profile{}, &sharedState)
	assert.Equal(t, 0, es.Len())
	assert.Panics(t, func() { es.AppendEmpty() })
	assert.Panics(t, func() { es.EnsureCapacity(2) })
	es2 := NewProfilesSlice()
	es.CopyTo(es2)
	assert.Panics(t, func() { es2.CopyTo(es) })
	assert.Panics(t, func() { es.MoveAndAppendTo(es2) })
	assert.Panics(t, func() { es2.MoveAndAppendTo(es) })
}

func TestProfilesSlice_CopyTo(t *testing.T) {
	dest := NewProfilesSlice()
	// Test CopyTo to empty
	NewProfilesSlice().CopyTo(dest)
	assert.Equal(t, NewProfilesSlice(), dest)

	// Test CopyTo larger slice
	generateTestProfilesSlice().CopyTo(dest)
	assert.Equal(t, generateTestProfilesSlice(), dest)

	// Test CopyTo same size slice
	generateTestProfilesSlice().CopyTo(dest)
	assert.Equal(t, generateTestProfilesSlice(), dest)
}

func TestProfilesSlice_EnsureCapacity(t *testing.T) {
	es := generateTestProfilesSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestProfilesSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestProfilesSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestProfilesSlice(), es)
}

func TestProfilesSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestProfilesSlice()
	dest := NewProfilesSlice()
	src := generateTestProfilesSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestProfilesSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestProfilesSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestProfilesSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestProfilesSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewProfilesSlice()
	emptySlice.RemoveIf(func(el Profile) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestProfilesSlice()
	pos := 0
	filtered.RemoveIf(func(el Profile) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestProfilesSliceAll(t *testing.T) {
	ms := generateTestProfilesSlice()
	assert.NotEmpty(t, ms.Len())

	var c int
	for i, v := range ms.All() {
		assert.Equal(t, ms.At(i), v, "element should match")
		c++
	}
	assert.Equal(t, ms.Len(), c, "All elements should have been visited")
}

func TestProfilesSlice_Sort(t *testing.T) {
	es := generateTestProfilesSlice()
	es.Sort(func(a, b Profile) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.Less(t, uintptr(unsafe.Pointer(es.At(i-1).orig)), uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b Profile) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.Greater(t, uintptr(unsafe.Pointer(es.At(i-1).orig)), uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}

func generateTestProfilesSlice() ProfilesSlice {
	es := NewProfilesSlice()
	fillTestProfilesSlice(es)
	return es
}

func fillTestProfilesSlice(es ProfilesSlice) {
	*es.orig = make([]*otlpprofiles.Profile, 7)
	for i := 0; i < 7; i++ {
		(*es.orig)[i] = &otlpprofiles.Profile{}
		fillTestProfile(newProfile((*es.orig)[i], es.state))
	}
}
