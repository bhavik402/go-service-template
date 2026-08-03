package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"{{ module_name }}/internal/domain"
	"{{ module_name }}/internal/messaging/events"
)

func (s *UserService) publishUserCreated(ctx context.Context, user *domain.User) {
	payload, err := json.Marshal(events.UserCreated{
		UserID:     user.ID.String(),
		Email:      user.Email,
		OccurredAt: time.Now().UTC(),
	})
	if err != nil {
		s.logger.Warn("failed to marshal user.created event", slog.Any("error", err))
		return
	}

	if err := s.producer.Publish(ctx, events.TopicUserCreated, user.ID.String(), payload); err != nil {
		s.logger.Warn("failed to publish user.created event",
			slog.Any("error", err), slog.String("user_id", user.ID.String()))
	}
}
