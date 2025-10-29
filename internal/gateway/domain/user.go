package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/users"
)

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      users.Role
	Nickname  string
	Email     string
	PhotoURL  string
}
