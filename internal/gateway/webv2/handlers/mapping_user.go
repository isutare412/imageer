package handlers

import (
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
)

func UserToWeb(p domain.User) gen.User {
	return gen.User{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Role:      p.Role,
		Nickname:  p.Nickname,
		Email:     p.Email,
		PhotoURL:  p.PhotoURL,
	}
}
