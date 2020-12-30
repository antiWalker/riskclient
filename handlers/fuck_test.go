package handlers

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFuck(t *testing.T) {
	a := FieldForm{Field: "hello"}

	b := func(v interface{}) {

		if data, err := json.Marshal(v); err == nil {
			t.Log(string(data))
			typ := reflect.TypeOf(v)
			q := reflect.New(typ).Interface()

			if err := json.Unmarshal(data, q); err == nil {
				t.Log(a)
				t.Log(q)
			}
		}
	}

	b(a)
}

func p(value reflect.Value, t *testing.T) {
	switch value.Kind() {
	case reflect.Bool:
		t.Log("bool")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		t.Log("int")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		t.Log("uint")
	case reflect.Uintptr:
		t.Log("uintptr")
	case reflect.Float32, reflect.Float64:
		t.Log("float")
	case reflect.Complex64, reflect.Complex128:
		t.Log("complex")
	case reflect.Array:
		t.Log("array")
	case reflect.Chan:
		t.Log("chan")
	case reflect.Func:
		t.Log("func")
	case reflect.Interface:
		t.Log("interface")
		p(value.Elem(), t)
	case reflect.Map:
		t.Log("map")
		for _, key := range value.MapKeys() {
			v := value.MapIndex(key)
			p(v, t)
		}
	case reflect.Ptr:
		t.Log("ptr")
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			p(value.Index(i), t)
		}
		t.Log("slice")
	case reflect.String:
		t.Log("string")
	case reflect.Struct:
		t.Log("struct")
	case reflect.UnsafePointer:
		t.Log("pointer")
	default:
		t.Log("fuck")
	}
}

func Test_JsonHello(t *testing.T) {
	a := "{\"hello\": [1, 2, 3], \"world\": {\"1\": 2}}"

	b := map[string]interface{}{}

	_ = json.Unmarshal([]byte(a), &b)

	t.Log(b)

	value := reflect.ValueOf(b)

	p(value, t)
}
