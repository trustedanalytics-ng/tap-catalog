package models

import "reflect"

func init() {
	RegisterType("Plans", reflect.TypeOf(ServicePlan{}))
	RegisterType("Cost", reflect.TypeOf(ServicePlanCost{}))
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
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Cost        ServicePlanCost `json:"cost"`
	AuditTrail  AuditTrail      `json:"auditTrail"`
}

type ServicePlanCost struct {
	Currency string `json:"currency"`
	//TODO DPNG-8533 define other attributes of cost
}

type ServiceState string

const (
	ServiceStateDeploying	ServiceState = "DEPLOYING"
	ServiceStateReady	ServiceState = "READY"
	ServiceStateOffline	ServiceState = "OFFLINE"
)
