package models

import (
	"encoding/json"
	"reflect"
)

var Registry = map[string]reflect.Type{}

func RegisterType(name string, t reflect.Type) {
	Registry[name] = t
}

type Patch struct {
	Operation string          `json:"op"`
	Field     string          `json:"field"`
	Value     json.RawMessage `json:"value"`
}
