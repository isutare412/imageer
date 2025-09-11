package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/users"
)

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Authority users.Authority
	Nickname  string
	Email     string
	PhotoURL  string
}
