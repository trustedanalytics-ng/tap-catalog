package models

import "reflect"

func init() {
	RegisterType("Metadata", reflect.TypeOf(Metadata{}))
	RegisterType("Bindings", reflect.TypeOf(InstanceBindings{}))
}

type Instance struct {
	Id         string             `json:"id"`
	Name       string             `json:"name"`
	Type       InstanceType       `json:"type"`
	ClassId    string             `json:"classId"`
	Bindings   []InstanceBindings `json:"bindings"`
	Metadata   []Metadata         `json:"metadata"`
	State      InstanceState      `json:"state"`
	AuditTrail AuditTrail         `json:"auditTrail"`
}

type InstanceState string

const (
	InstanceStateRequested   InstanceState = "REQUESTED"
	InstanceStateDeploying   InstanceState = "DEPLOYING"
	InstanceStateFailure     InstanceState = "FAILURE"
	InstanceStateStopped     InstanceState = "STOPPED"
	InstanceStateStartReq    InstanceState = "START_REQ"
	InstanceStateStarting    InstanceState = "STARTING"
	InstanceStateRunning     InstanceState = "RUNNING"
	InstanceStateStopReq     InstanceState = "STOP_REQ"
	InstanceStateStopping    InstanceState = "STOPPING"
	InstanceStateDestroyReq  InstanceState = "DESTROY_REQ"
	InstanceStateDestroying  InstanceState = "DESTROYING"
	InstanceStateUnavailable InstanceState = "UNAVAILABLE"
)

type InstanceBindings struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}

type Metadata struct {
	Id    string `json:"key"`
	Value string `json:"value"`
}

type InstanceType string

const (
	InstanceTypeApplication   InstanceType = "APPLICATION"
	InstanceTypeService       InstanceType = "SERVICE"
	InstanceTypeServiceBroker InstanceType = "SERVICE_BROKER"
)
