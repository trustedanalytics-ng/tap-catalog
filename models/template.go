package models

type Template struct {
	Id         string     `json:"templateId"`
	State      string     `json:"state"`
	AuditTrail AuditTrail `json:"auditTrail"`
}

//todo add State type
