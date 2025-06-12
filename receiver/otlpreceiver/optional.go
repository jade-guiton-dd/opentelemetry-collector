package otlpreceiver

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/confmap"
)

type Optional[T any] struct {
	Enabled bool
	Value   T
}

func NoneWithDefault[T any](val T) Optional[T] {
	return Optional[T]{
		Enabled: false,
		Value:   val,
	}
}
func Some[T any](val T) Optional[T] {
	return Optional[T]{
		Enabled: true,
		Value:   val,
	}
}

var _ confmap.FieldUnmarshaler = (*Optional[configgrpc.ServerConfig])(nil)

func (o *Optional[T]) UnmarshalField(data any, isSet bool, meta confmap.FieldMetadata) error {
	if isSet {
		if dataMap, ok := data.(map[string]any); ok {
			if err := confmap.NewFromStringMap(dataMap).Unmarshal(&o.Value); err != nil {
				return err
			}
		} else if data == nil {
			// treat null like empty map, ie. keep enable and keep defaults
		} else {
			return fmt.Errorf("'protocols::%s' expected a map or null, got '%s'", meta.Key, reflect.ValueOf(data).Kind())
		}
		o.Enabled = true
	} else {
		o.Enabled = false
	}
	return nil
}
