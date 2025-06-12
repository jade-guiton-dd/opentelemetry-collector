package confmap

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

// The FieldUnmarshaler interface allows config structs to bypass mapstructure's
// unmarshalling behavior and run arbitrary code when a type is used as a field
// of a struct being unmarshalled.
//
// The `UnmarshalField` method will be called:
//   - if the field is not set in the input config, in which case `isSet == false`
//   - if the field is set to a value that would normally cause an error, such as
//     an integer for a struct field
//
// This allows implementing custom unmarshalling behaviors, such as providing
// default values for subfields of an optional field.
type FieldUnmarshaler interface {
	UnmarshalField(data any, isSet bool, meta FieldMetadata) error
}

type FieldMetadata struct {
	Key string
}

func fieldUnmarshalerHookFunc() mapstructure.DecodeHookFuncValue {
	return safeWrapDecodeHookFunc(func(from reflect.Value, to reflect.Value) (any, error) {
		ty := to.Type()
		if ty.Implements(reflect.TypeFor[FieldUnmarshaler]()) {
			return nil, fmt.Errorf("FieldUnmarshaler type '%s' was used outside of a struct field", ty.Name())
		}
		if ty.Kind() != reflect.Struct {
			return from.Interface(), nil
		}
		fromAsMap, ok := from.Interface().(map[string]any)
		if !ok {
			return from.Interface(), nil
		}
		for i := 0; i < to.Type().NumField(); i++ {
			fieldType := to.Type().Field(i)
			if !fieldType.IsExported() {
				continue
			}
			if unmarshaler, ok := to.Field(i).Addr().Interface().(FieldUnmarshaler); ok {
				key := fieldType.Name
				tag := fieldType.Tag.Get("mapstructure")
				tag, _, _ = strings.Cut(tag, ",")
				if tag != "" {
					key = tag
				}

				val, isSet := fromAsMap[key]
				if isSet {
					delete(fromAsMap, key) // so mapstructure doesn't touch the result
				}
				meta := FieldMetadata{
					Key: key,
				}
				if err := unmarshaler.UnmarshalField(val, isSet, meta); err != nil {
					return nil, err
				}
			}
		}
		return fromAsMap, nil
	})
}
