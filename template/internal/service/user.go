// Package service holds business logic and orchestration. Services
// depend only on interfaces (repository, client), never on concrete
// infrastructure — this file, cache-aside reads and event publishing
// aside, has no idea whether persistence is Postgres or something else.
package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"{{ module_name }}/internal/client/notification"
	"{{ module_name }}/internal/domain"
	"{{ module_name }}/internal/repository"
{% if use_redis %}
	"{{ module_name }}/internal/cache"
{% endif %}
{% if use_kafka %}
	"{{ module_name }}/internal/messaging/kafka"
{% endif %}
)

type UserService struct {
	repo     repository.UserRepository
	notifier notification.Notifier
{% if use_redis %}
	cache cache.Cache
{% endif %}
{% if use_kafka %}
	producer *kafka.Producer
{% endif %}
	logger *slog.Logger
}

func NewUserService(
	repo repository.UserRepository,
	notifier notification.Notifier,
{% if use_redis %}
	cache cache.Cache,
{% endif %}
{% if use_kafka %}
	producer *kafka.Producer,
{% endif %}
	logger *slog.Logger,
) *UserService {
	return &UserService{
		repo:     repo,
		notifier: notifier,
{% if use_redis %}
		cache: cache,
{% endif %}
{% if use_kafka %}
		producer: producer,
{% endif %}
		logger: logger,
	}
}

type CreateUserInput struct {
	Name  string
	Email string
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	now := time.Now().UTC()
	user := &domain.User{
		ID:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
{% if use_kafka %}

	s.publishUserCreated(ctx, user)
{% endif %}

	// Best-effort: a failed welcome email should never fail user creation.
	if _, err := s.notifier.SendWelcomeEmail(ctx, notification.SendWelcomeEmailRequest{
		ToEmail: user.Email,
		ToName:  user.Name,
	}); err != nil {
		s.logger.Warn("failed to send welcome email",
			slog.Any("error", err), slog.String("user_id", user.ID.String()))
	}

	return user, nil
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
{% if use_redis %}
	if cached, err := s.getCached(ctx, id); err == nil {
		return cached, nil
	}
{% endif %}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
{% if use_redis %}

	s.setCached(ctx, user)
{% endif %}

	return user, nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	return s.repo.List(ctx, limit, offset)
}

type UpdateUserInput struct {
	Name  string
	Email string
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Name = input.Name
	user.Email = input.Email
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
{% if use_redis %}

	s.invalidateCache(ctx, id)
{% endif %}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
{% if use_redis %}

	s.invalidateCache(ctx, id)
{% endif %}

	return nil
}
