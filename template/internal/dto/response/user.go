package response

import (
	"time"

	"github.com/google/uuid"

	"{{ module_name }}/internal/domain"
	"{{ module_name }}/internal/util"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromDomainUser(u *domain.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserList struct {
	Items []User        `json:"items"`
	Meta  util.PageMeta `json:"meta"`
}
