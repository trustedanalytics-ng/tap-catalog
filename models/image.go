package models

import "strings"

const USER_DEFINED_APPLICATION_IMAGE_PREFIX = "app_"
const USER_DEFINED_OFFERING_IMAGE_PREFIX = "svc_"

type Image struct {
	Id         string     `json:"id"`
	Type       ImageType  `json:"type"`
	BlobType   BlobType   `json:"blobType"`
	State      ImageState `json:"state"`
	AuditTrail AuditTrail `json:"auditTrail"`
}

type ImageType string

const (
	ImageTypeJava     ImageType = "JAVA"
	ImageTypeGo       ImageType = "GO"
	ImageTypeNodeJs   ImageType = "NODEJS"
	ImageTypePython27 ImageType = "PYTHON2.7"
	ImageTypePython34 ImageType = "PYTHON3.4"
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
	return strings.HasPrefix(imageId, USER_DEFINED_APPLICATION_IMAGE_PREFIX)
}

func GetApplicationId(imageId string) string {
	return strings.TrimPrefix(imageId, USER_DEFINED_APPLICATION_IMAGE_PREFIX)
}

func GenerateImageId(applicationId string) string {
	return USER_DEFINED_APPLICATION_IMAGE_PREFIX + applicationId
}

func IsUserDefinedOffering(imageId string) bool {
	return strings.HasPrefix(imageId, USER_DEFINED_OFFERING_IMAGE_PREFIX)
}

func GetOfferingId(imageId string) string {
	return strings.TrimPrefix(imageId, USER_DEFINED_OFFERING_IMAGE_PREFIX)
}

func ConstructImageIdForUserOffering(offeringId string) string {
	return USER_DEFINED_OFFERING_IMAGE_PREFIX + offeringId
}