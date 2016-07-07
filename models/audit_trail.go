package models

type AuditTrail struct {
	//todo change int to time
	CreatedOn     int    `json:"createdOn"`
	CreatedBy     string `json:"createdBy"`
	LastUpdatedOn int    `json:"lastUpdatedOn"`
}
