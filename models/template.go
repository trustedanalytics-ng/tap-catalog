package models

type Template struct {
	Id         string        `json:"templateId"`
	State      TemplateState `json:"state"`
	AuditTrail AuditTrail    `json:"auditTrail"`
}

type TemplateState string

const (
	TemplateStateInProgress  TemplateState = "IN_PROGRESS"
	TemplateStateReady       TemplateState = "READY"
	TemplateStateUnavailable TemplateState = "UNAVAILABLE"
)
