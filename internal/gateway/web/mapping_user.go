package web

import "github.com/isutare412/imageer/internal/gateway/domain"

func UserToWeb(p domain.User) User {
	return User{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Role:      p.Role,
		Nickname:  p.Nickname,
		Email:     p.Email,
		PhotoURL:  p.PhotoURL,
	}
}
