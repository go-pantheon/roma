package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const (
	SliceSep = "#"
	KvSep    = "$"
)

type Field struct {
	*Metadata

	RawValue string
	Value    any
}

func NewField(md *Metadata, raw string) (fd *Field, err error) {
	v, err := stringToValue(raw, md.FieldType)
	if err != nil {
		err = errors.Wrapf(err, "field=%s raw=%s parse value failed", md.FieldName, raw)
		return
	}

	fd = &Field{
		Metadata: md,
		RawValue: raw,
		Value:    v,
	}
	return fd, nil
}

func (f *Field) EncodeToJson() (string, error) {
	name := f.FieldName
	if f.Value == nil {
		return "", errors.Errorf("field value must not nil. field=%s", f.FieldName)
	}

	jsonData, err := json.Marshal(f.Value)
	if err != nil {
		return "", errors.Wrapf(err, "failed to encode field to JSON. field=%s", f.FieldName)
	}
	return fmt.Sprintf("\"%s\": %s", name, string(jsonData)), nil
}

func toReflectType(typeStr string) (reflect.Type, error) {
	switch typeStr {
	case "int":
		return reflect.TypeOf(int64(0)), nil
	case "[]int":
		return reflect.TypeOf([]int64{}), nil
	case "uint":
		return reflect.TypeOf(uint64(0)), nil
	case "[]uint":
		return reflect.TypeOf([]uint64{}), nil
	case "string":
		return reflect.TypeOf(""), nil
	case "[]string":
		return reflect.TypeOf([]string{}), nil
	case "bool":
		return reflect.TypeOf(true), nil
	case "[]bool":
		return reflect.TypeOf([]bool{}), nil
	case "float":
		return reflect.TypeOf(0.0), nil
	case "[]float":
		return reflect.TypeOf([]float64{}), nil
	case "map[int]int":
		return reflect.TypeOf(map[int64]int64{}), nil
	case "map[int]uint":
		return reflect.TypeOf(map[int64]uint64{}), nil
	case "map[uint]int":
		return reflect.TypeOf(map[uint64]int64{}), nil
	case "map[uint]uint":
		return reflect.TypeOf(map[uint64]uint64{}), nil
	case "map[uint]string":
		return reflect.TypeOf(map[uint64]string{}), nil
	case "map[string]int":
		return reflect.TypeOf(map[uint64]int64{}), nil
	case "map[string]uint":
		return reflect.TypeOf(map[uint64]uint64{}), nil
	case "map[string]string":
		return reflect.TypeOf(map[uint64]string{}), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeStr)
	}
}

func stringToValue(str string, t reflect.Type) (any, error) {
	switch t.Kind() {
	case reflect.Int64:
		var v int64
		if str == "" {
			return int64(0), nil
		}
		_, err := fmt.Sscanf(str, "%d", &v)
		return v, err
	case reflect.Uint64:
		var v uint64
		if str == "" {
			return uint64(0), nil
		}
		_, err := fmt.Sscanf(str, "%d", &v)
		return v, err
	case reflect.String:
		return str, nil
	case reflect.Bool:
		var v bool
		if str == "" {
			return false, nil
		}
		_, err := fmt.Sscanf(str, "%t", &v)
		return v, err
	case reflect.Float64:
		var v float64
		if str == "" {
			return float64(0), nil
		}
		_, err := fmt.Sscanf(str, "%f", &v)
		return v, err
	case reflect.Slice:
		if str == "" {
			return reflect.MakeSlice(t, 0, 0).Interface(), nil
		}
		elemType := t.Elem()
		elements := strings.Split(str, SliceSep)
		slice := reflect.MakeSlice(t, len(elements), len(elements))
		for i, elemStr := range elements {
			elemValue, err := stringToValue(elemStr, elemType)
			if err != nil {
				return nil, err
			}
			slice.Index(i).Set(reflect.ValueOf(elemValue))
		}
		return slice.Interface(), nil
	case reflect.Map:
		if str == "" {
			return reflect.MakeMap(t).Interface(), nil
		}
		mapType := t
		keyType := mapType.Key()
		elemType := mapType.Elem()
		mapValue := reflect.MakeMap(mapType)
		pairs := strings.Split(str, SliceSep)
		for _, pair := range pairs {
			kv := strings.Split(pair, KvSep)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid key-value pair<%s>", pair)
			}
			key, err := stringToValue(kv[0], keyType)
			if err != nil {
				return nil, err
			}
			value, err := stringToValue(kv[1], elemType)
			if err != nil {
				return nil, err
			}
			if v := mapValue.MapIndex(reflect.ValueOf(key)); v.IsValid() {
				return nil, fmt.Errorf("duplicate map key<%v> value<%v>", key, v)
			}
			mapValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
		}
		return mapValue.Interface(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", t)
	}
}

func IsZeroValue(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	default:
		return reflect.ValueOf(v).IsZero()
	}
}
