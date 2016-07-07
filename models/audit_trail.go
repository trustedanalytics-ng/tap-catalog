package models

type AuditTrail struct {
	//todo change int to time
	Id            string `json:"id"`
	CreatedOn     int    `json:"createdOn"`
	CreatedBy     string `json:"createdBy"`
	LastUpdatedOn int    `json:"lastUpdatedOn"`
}
