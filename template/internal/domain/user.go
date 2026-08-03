// Package domain holds the core business models. It has no framework or
// infrastructure imports — no HTTP, no SQL, no JSON tags driven by the
// wire format. dto/ and repository/ translate to and from these types.
package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
