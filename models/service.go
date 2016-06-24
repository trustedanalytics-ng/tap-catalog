package models

type Service struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Bindable    bool          `json:"bindable"`
	TemplateId  string        `json:"templateId"`
	State       string        `json:"state"`
	Plans       []ServicePlan `json:"plans"`
	AuditTrail  AuditTrail
}

type ServicePlan struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Cost        ServicePlanCost `json:"cost"`
	AuditTrail  AuditTrail
}

type ServicePlanCost struct {
	Currency string `json:"currency"`
	//TODO define other attributes of cost
}
