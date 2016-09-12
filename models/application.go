package models

type Application struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ImageId     string     `json:"imageId"`
	Replication int        `json:"replication"`
	TemplateId  string     `json:"templateId"`
	AuditTrail  AuditTrail `json:"auditTrail"`
}
