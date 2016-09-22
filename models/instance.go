package models

import (
	"reflect"
	"strings"
)

func init() {
	RegisterType("Metadata", reflect.TypeOf(Metadata{}))
	RegisterType("Bindings", reflect.TypeOf(InstanceBindings{}))
}

const (
	BROKER_OFFERING_PREFIX    = "BROKER_OFFERING_"
	BROKER_TEMPLATE_ID        = "BROKER_TEMPLATE_ID"
	APPLICATION_IMAGE_ADDRESS = "APPLICATION_IMAGE_ADDRESS"
	OFFERING_PLAN_ID          = "PLAN_ID"
)

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

func GetValueFromMetadata(metadatas []Metadata, key string) string {
	for _, metadata := range metadatas {
		if metadata.Id == key {
			return metadata.Value
		}
	}
	return ""
}

func GetPrefixedOfferingName(offeringName string) string {
	return BROKER_OFFERING_PREFIX + offeringName
}

func IsServiceBrokerOfferingMetadata(metadata Metadata) bool {
	return strings.Contains(metadata.Id, BROKER_OFFERING_PREFIX)
}
