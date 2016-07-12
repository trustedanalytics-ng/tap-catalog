package models

type AuditTrail struct {
	CreatedOn     int64  `json:"createdOn"`
	CreatedBy     string `json:"createdBy"`
	LastUpdatedOn int64  `json:"lastUpdatedOn"`
	LastUpdateBy  string `json:"lastUpdateBy"`
}
