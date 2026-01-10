package entity

import "time"

type ImageProcessingLog struct {
	ID             int
	CreatedAt      time.Time
	IsSuccess      bool
	ErrorCode      int     `gorm:"type:integer"`
	ErrorMessage   *string `gorm:"type:text"`
	DurationMillis int     `gorm:"type:integer"`

	ImageVariantID string `gorm:"size:36"`
}
