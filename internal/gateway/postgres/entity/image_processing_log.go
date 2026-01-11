package entity

import (
	"time"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

type ImageProcessingLog struct {
	ID                int
	CreatedAt         time.Time
	IsSuccess         bool
	ErrorCode         *int    `gorm:"type:integer"`
	ErrorMessage      *string `gorm:"type:text"`
	ElapsedTimeMillis int     `gorm:"type:integer"`
	ImageVariantID    string  `gorm:"size:36"`
}

func NewImageProcessingLog(log domain.ImageProcessingLog) ImageProcessingLog {
	return ImageProcessingLog{
		IsSuccess:         log.IsSuccess,
		ErrorCode:         log.ErrorCode,
		ErrorMessage:      log.ErrorMessage,
		ElapsedTimeMillis: int(log.ElapsedTime.Milliseconds()),
		ImageVariantID:    log.ImageVariantID,
	}
}

func (l ImageProcessingLog) ToDomain() domain.ImageProcessingLog {
	return domain.ImageProcessingLog{
		ID:             l.ID,
		CreatedAt:      l.CreatedAt,
		IsSuccess:      l.IsSuccess,
		ErrorCode:      l.ErrorCode,
		ErrorMessage:   l.ErrorMessage,
		ElapsedTime:    time.Duration(l.ElapsedTimeMillis) * time.Millisecond,
		ImageVariantID: l.ImageVariantID,
	}
}
