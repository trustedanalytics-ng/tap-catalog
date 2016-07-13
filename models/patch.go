package models

import (
	"encoding/json"
	"reflect"
)

var Registry = map[string]reflect.Type{}

func RegisterType(name string, t reflect.Type) {
	Registry[name] = t
}

type PatchOperation string

const (
	OperationAdd    PatchOperation = "Add"
	OperationUpdate PatchOperation = "Update"
	OperationDelete PatchOperation = "Delete"
)

type Patch struct {
	Operation PatchOperation  `json:"op"`
	Field     string          `json:"field"`
	Value     json.RawMessage `json:"value"`
	PrevValue json.RawMessage `json:"prevValue,omitempty"`
}
