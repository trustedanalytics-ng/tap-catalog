package models

type Instance struct {
	Id         string             `json:"id"`
	Type       InstanceType       `json:"type"`
	ClassId    string             `json:"classId"`
	Bindings   []InstanceBindings `json:"bindings"`
	Metadata   []InstanceMetadata `json:"meta"`
	State      InstanceState      `json:"state"`
	AuditTrail AuditTrail
}

type InstanceState string

const (
	InstanceStateRequested     InstanceState = "requested"
	InstanceStateDeploying     InstanceState = "deploying"
	InstanceStateFailure       InstanceState = "failure"
	InstanceStateStopped       InstanceState = "stopped"
	InstanceStateRunning       InstanceState = "running"
	InstanceStateToBeDestroyed InstanceState = "tobedestroyed"
	InstanceStateDestroying    InstanceState = "destroying"
	InstanceStateUnavailable   InstanceState = "unavailable"
)

type InstanceBindings struct {
	Id string `json:"id"`
}

type InstanceMetadata struct {
	Id    string `json:"key"`
	Value string `json:"value"`
}

type InstanceType string

const (
	InstanceTypeApplication   InstanceType = "application"
	InstanceTypeService       InstanceType = "service"
	InstanceTypeServiceBroker InstanceType = "service_broker"
)
