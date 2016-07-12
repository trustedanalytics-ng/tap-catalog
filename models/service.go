package models

import "reflect"

func init() {
	RegisterType("Plans", reflect.TypeOf(ServicePlan{}))
}

type Service struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Bindable    bool          `json:"bindable"`
	TemplateId  string        `json:"templateId"`
	State       ServiceState  `json:"state"`
	Plans       []ServicePlan `json:"plans"`
	AuditTrail  AuditTrail    `json:"auditTrail"`
}

type ServicePlan struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Cost        string     `json:"cost"`
	AuditTrail  AuditTrail `json:"auditTrail"`
}

type ServiceState string

const (
	ServiceStateDeploying ServiceState = "DEPLOYING"
	ServiceStateReady     ServiceState = "READY"
	ServiceStateOffline   ServiceState = "OFFLINE"
)
