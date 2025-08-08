package pcommon

import (
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

func (m Map) PutEmptyUnsafe(k string) Value {
	m.getState().AssertMutable()
	*m.getOrig() = append(*m.getOrig(), otlpcommon.KeyValue{Key: k})
	return newValue(&(*m.getOrig())[len(*m.getOrig())-1].Value, m.getState())
}
