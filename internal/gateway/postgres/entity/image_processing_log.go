package entity

import "time"

type ImageProcessingLog struct {
	ID             int
	CreatedAt      time.Time
	IsSuccess      bool
	ErrorMessage   *string `gorm:"type:text"`
	DurationMillis int     `gorm:"type:integer"`

	ImageVariantID string `gorm:"size:36"`
}
