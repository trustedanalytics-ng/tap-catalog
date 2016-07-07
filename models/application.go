package models

type Application struct {
	Id          string     `json:"id"`
	ImageId     string     `json:"imageId"`
	Replication int        `json:"replication"`
	TemplateId  string     `json:"templateId"`
	AuditTrail  AuditTrail `json:"auditTrail"`
}
