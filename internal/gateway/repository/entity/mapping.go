package entity

import "github.com/isutare412/imageer/internal/gateway/domain"

func UserToDomain(u *User) domain.User {
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
