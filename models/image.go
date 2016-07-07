package models

type Image struct {
	Id         string     `json:"id"`
	Type       ImageType  `json:"type"`
	State      ImageState `json:"state"`
	AuditTrail AuditTrail `json:"auditTrail"`
}

type ImageType string

const (
	ImageTypeJava   ImageType = "JAVA"
	ImageTypeGo     ImageType = "GO"
	ImageTypeNodeJs ImageType = "NODEJS"
)

type ImageState string

const (
	ImageStatePending ImageState = "PENDING"
	ImageTypeBuilding ImageState = "BUILDING"
	ImageTypeError    ImageState = "ERROR"
	ImageTypeReady    ImageState = "READY"
)
