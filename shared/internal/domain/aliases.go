package domain

import (
	infra "github.com/brunojet/infra-go/pkg/domain"
)

const (
	ObjectStatusPending    = infra.ObjectStatusPending
	ObjectStatusProcessing = infra.ObjectStatusProcessing
	ObjectStatusAvailable  = infra.ObjectStatusAvailable
	ObjectStatusError      = infra.ObjectStatusError
	ImageTypeIcon          = infra.ImageTypeIcon
	ImageTypeScreenshot    = infra.ImageTypeScreenshot
	ImageTypeBanner        = infra.ImageTypeBanner
	VideoTypeShorts        = infra.VideoTypeShorts
	VideoTypeReels         = infra.VideoTypeReels
)

type ObjectStatus = infra.ObjectStatus
type VideoType = infra.VideoType
type ImageType = infra.ImageType
type BaseModel = infra.BaseModel
type BaseEntity = infra.BaseEntity
type StorageObject = infra.StorageObject
type Video = infra.Video
type Image = infra.Image
