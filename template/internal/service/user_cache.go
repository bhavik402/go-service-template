package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"{{ module_name }}/internal/domain"
)

const userCacheTTL = 5 * time.Minute

func userCacheKey(id uuid.UUID) string {
	return "user:" + id.String()
}

func (s *UserService) getCached(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	raw, err := s.cache.Get(ctx, userCacheKey(id))
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(raw), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) setCached(ctx context.Context, user *domain.User) {
	raw, err := json.Marshal(user)
	if err != nil {
		return
	}
	if err := s.cache.Set(ctx, userCacheKey(user.ID), string(raw), userCacheTTL); err != nil {
		s.logger.Warn("failed to cache user", slog.Any("error", err), slog.String("user_id", user.ID.String()))
	}
}

func (s *UserService) invalidateCache(ctx context.Context, id uuid.UUID) {
	if err := s.cache.Delete(ctx, userCacheKey(id)); err != nil {
		s.logger.Warn("failed to invalidate user cache", slog.Any("error", err), slog.String("user_id", id.String()))
	}
}
