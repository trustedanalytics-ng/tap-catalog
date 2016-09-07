package models

import "reflect"

func init() {
	RegisterType("Plans", reflect.TypeOf(ServicePlan{}))
}

const (
	ENV_SOURCE_OFFERING_ID       = "source_offering_id"
	ENV_BROKER_SHORT_INSTANCE_ID = "broker_short_instance_id"
	ENV_BROKER_INSTANCE_ID       = "broker_instance_id"
	ENV_PLAN_ID                  = "plan_id"
	ENV_SOURCE_PLAN_ID_PREFIX    = "source_plan_id-"
)

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

func GetPrefixedSourcePlanName(planName string) string {
	return ENV_SOURCE_PLAN_ID_PREFIX + planName
}
