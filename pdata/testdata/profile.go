// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package testdata // import "go.opentelemetry.io/collector/pdata/testdata"

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pprofile"
)

var (
	profileStartTimestamp = pcommon.NewTimestampFromTime(time.Date(2020, 2, 11, 20, 26, 12, 321, time.UTC))
	profileEndTimestamp   = pcommon.NewTimestampFromTime(time.Date(2020, 2, 11, 20, 26, 13, 789, time.UTC))
)

// GenerateProfiles generates dummy profiling data for tests
func GenerateProfiles(profilesCount int) pprofile.Profiles {
	td := pprofile.NewProfiles()
	initResource(td.ResourceProfiles().AppendEmpty().Resource())
	ss := td.ResourceProfiles().At(0).ScopeProfiles().AppendEmpty().Profiles()

	dic := td.ProfilesDictionary()
	dic.StringTable().Append("")
	attr := dic.AttributeTable().AppendEmpty()
	attr.SetKey("key")
	attr.Value().SetStr("value")
	attr2 := dic.AttributeTable().AppendEmpty()
	attr2.SetKey("key")
	attr2.Value().SetStr("value")

	ss.EnsureCapacity(profilesCount)
	for i := 0; i < profilesCount; i++ {
		switch i % 2 {
		case 0:
			fillProfileOne(dic, ss.AppendEmpty())
		case 1:
			fillProfileTwo(dic, ss.AppendEmpty())
		}
	}

	return td
}

func fillProfileOne(dic pprofile.ProfilesDictionary, profile pprofile.Profile) {
	profile.SetProfileID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	profile.SetTime(profileStartTimestamp)
	profile.SetDuration(profileEndTimestamp)
	profile.SetDroppedAttributesCount(1)
	profile.PeriodType().SetAggregationTemporality(pprofile.AggregationTemporalityDelta)

	st := profile.SampleType().AppendEmpty()
	st.SetAggregationTemporality(pprofile.AggregationTemporalityDelta)

	loc := pprofile.NewLocation()
	loc.SetAddress(1)
	_ = pprofile.PutLocation(dic.LocationTable(), profile, loc)

	sample := profile.Sample().AppendEmpty()
	sample.SetLocationsStartIndex(0)
	sample.SetLocationsLength(1)
	sample.Value().Append(4)
	sample.AttributeIndices().Append(0)
}

func fillProfileTwo(dic pprofile.ProfilesDictionary, profile pprofile.Profile) {
	profile.SetProfileID([16]byte{0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	profile.SetTime(profileStartTimestamp)
	profile.SetDuration(profileEndTimestamp)
	profile.PeriodType().SetAggregationTemporality(pprofile.AggregationTemporalityDelta)

	st := profile.SampleType().AppendEmpty()
	st.SetAggregationTemporality(pprofile.AggregationTemporalityDelta)

	loc := pprofile.NewLocation()
	loc.SetAddress(2)
	_ = pprofile.PutLocation(dic.LocationTable(), profile, loc)

	sample := profile.Sample().AppendEmpty()
	sample.SetLocationsStartIndex(0)
	sample.SetLocationsLength(1)
	sample.Value().Append(9)
	sample.AttributeIndices().Append(0)
}
