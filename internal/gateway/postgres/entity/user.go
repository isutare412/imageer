package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/users"
)

type User struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      users.Role `gorm:"size:32"`
	Nickname  string     `gorm:"size:128"`
	Email     string     `gorm:"size:1024; uniqueIndex"`
	PhotoURL  string     `gorm:"size:2048"`
}

func NewUser(u domain.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Role:      u.Role,
		Nickname:  u.Nickname,
		Email:     u.Email,
		PhotoURL:  u.PhotoURL,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Role:      u.Role,
		Nickname:  u.Nickname,
		Email:     u.Email,
		PhotoURL:  u.PhotoURL,
	}
}
