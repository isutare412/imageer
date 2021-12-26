package user

import (
	"time"
)

type User struct {
	ID         int64  `gorm:"primaryKey"`
	GivenName  string `gorm:"not null; size:127"`
	FamilyName string `gorm:"size:127"`
	Email      string `gorm:"uniqueIndex; not null; size:127"`
	Password   string `gorm:"size:64"`
	Credit     int64  `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
