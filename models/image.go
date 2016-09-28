package models

import "strings"

const IMAGE_ID_PREFIX = "image-"

type Image struct {
	Id         string     `json:"id"`
	Type       ImageType  `json:"type"`
	BlobType   BlobType   `json:"blobType"`
	State      ImageState `json:"state"`
	AuditTrail AuditTrail `json:"auditTrail"`
}

type ImageType string

const (
	ImageTypeJava   ImageType = "JAVA"
	ImageTypeGo     ImageType = "GO"
	ImageTypeNodeJs ImageType = "NODEJS"
	ImageTypePython ImageType = "PYTHON"
)

type BlobType string

const (
	BlobTypeTarGz BlobType = "TARGZ"
	BlobTypeJar   BlobType = "JAR"
	BlobTypeExec  BlobType = "EXEC"
)

type ImageState string

const (
	ImageStatePending  ImageState = "PENDING"
	ImageStateBuilding ImageState = "BUILDING"
	ImageStateError    ImageState = "ERROR"
	ImageStateReady    ImageState = "READY"
)

func IsApplicationInstance(imageId string) bool {
	return strings.Contains(imageId, IMAGE_ID_PREFIX)
}

func GetApplicationId(imageId string) string {
	return strings.TrimPrefix(imageId, IMAGE_ID_PREFIX)
}

func GenerateImageId(applicationId string) string {
	return IMAGE_ID_PREFIX + applicationId
}
