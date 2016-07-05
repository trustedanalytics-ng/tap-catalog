package models

type Application struct {
	Id          string     `json:"id"`
	Image       string     `json:"image"`
	Replication int        `json:"replication"`
	Type        string     `json:"type"`
	TemplateId  string     `json:"templateId"`
	State       string     `json:"state"`
	AuditTrail  AuditTrail `json:"-"`
}
