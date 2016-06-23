package models

type Instance struct {
	Id       string             `json:"id"`
	Type     string             `json:"type"`
	ClassId  string             `json:"classId"`
	Bindings []InstanceBindings `json:"bindings"`
	Metadata []InstanceMetadata `json:"meta"`
	State    string             `json:"state"`
}

type InstanceBindings struct {
	Id string `json:"id"`
}

type InstanceMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
