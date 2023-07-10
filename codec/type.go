package codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

func TomlToJson(data []byte) ([]byte, error) {
	var val any
	if err := toml.NewDecoder(bytes.NewReader(data)).Decode(&val); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func YamlToJson(data []byte) ([]byte, error) {
	var val any
	if err := yaml.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	val = toStringKeyMap(val)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func convertKeyToString(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[Repr(k)] = toStringKeyMap(v)
	}
	return res
}

func convertNumberToJsonNumber(in any) json.Number {
	return json.Number(Repr(in))
}

func convertSlice(in []any) []any {
	res := make([]any, len(in))
	for i, v := range in {
		res[i] = toStringKeyMap(v)
	}
	return res
}

func toStringKeyMap(v any) any {
	switch v := v.(type) {
	case []any:
		return convertSlice(v)
	case map[any]any:
		return convertKeyToString(v)
	case bool, string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return convertNumberToJsonNumber(v)
	default:
		return Repr(v)
	}
}

func Repr(v interface{}) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return reprOfValue(val)
}

func reprOfValue(val reflect.Value) string {
	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}
