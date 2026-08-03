// Package events defines the payloads published to and consumed from
// messaging topics. These are wire contracts for other services, just
// like dto/ is the wire contract for HTTP clients — keep them separate
// from domain models so an internal refactor doesn't silently change an
// event schema.
package events

import "time"

const TopicUserCreated = "user.created"

type UserCreated struct {
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	OccurredAt time.Time `json:"occurred_at"`
}
