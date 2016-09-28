package models

import "reflect"

func init() {
	RegisterType("Plans", reflect.TypeOf(ServicePlan{}))
	RegisterType("Dependencies", reflect.TypeOf(ServiceDependency{}))
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
	Metadata    []Metadata    `json:"metadata"`
}

type ServicePlan struct {
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Cost         string              `json:"cost"`
	Dependencies []ServiceDependency `json:"dependencies"`
	AuditTrail   AuditTrail          `json:"auditTrail"`
}

type ServiceDependency struct {
	Id          string `json:"id"`
	PlanName    string `json:"plan_name"`
	PlanId      string `json:"plan_id"`
	ServiceName string `json:"service_name"`
	ServiceId   string `json:"service_id"`
}

type ServiceState string

const (
	ServiceStateDeploying ServiceState = "DEPLOYING"
	ServiceStateReady     ServiceState = "READY"
	ServiceStateOffline   ServiceState = "OFFLINE"
)
