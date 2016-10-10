package models

import "reflect"

func init() {
	RegisterType("InstanceDependencies", reflect.TypeOf(InstanceDependency{}))
}

type InstanceDependency struct {
	Id string `json:"id"`
}

type Application struct {
	Id                   string               `json:"id"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	ImageId              string               `json:"imageId"`
	Replication          int                  `json:"replication"`
	TemplateId           string               `json:"templateId"`
	AuditTrail           AuditTrail           `json:"auditTrail"`
	InstanceDependencies []InstanceDependency `json:"instanceDependencies"`
}
