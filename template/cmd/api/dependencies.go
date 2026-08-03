package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"{{ module_name }}/internal/client/notification"
	"{{ module_name }}/internal/config"
	"{{ module_name }}/internal/dto"
	"{{ module_name }}/internal/handler"
	"{{ module_name }}/internal/observability"
	"{{ module_name }}/internal/repository"
	"{{ module_name }}/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
{% if use_redis %}

	"{{ module_name }}/internal/cache"
	"{{ module_name }}/internal/locking"
{% endif %}
{% if use_kafka %}

	"{{ module_name }}/internal/messaging/events"
	"{{ module_name }}/internal/messaging/kafka"

	kafkago "github.com/segmentio/kafka-go"
{% endif %}
{% if use_s3 %}

	"{{ module_name }}/internal/storage"
{% endif %}
)

// Application is the fully wired object graph. cmd/api/main.go only calls
// BuildApplication and then Close — it never constructs a dependency on
// its own. This file is the composition root: the one place that knows
// how everything is wired together.
type Application struct {
	Config *config.Config
	Logger *slog.Logger
	Router http.Handler
{% if use_s3 %}
	Storage storage.Storage
{% endif %}
{% if use_redis %}
	Locker locking.Locker
{% endif %}

	pgPool *pgxpool.Pool
{% if use_redis %}
	redisCache *cache.RedisCache
{% endif %}
{% if use_kafka %}
	kafkaProducer  *kafka.Producer
	kafkaConsumer  *kafka.Consumer
	cancelConsumer context.CancelFunc
{% endif %}
}

func BuildApplication() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger := observability.NewLogger(cfg.LogLevel)
	metrics := observability.NewMetrics()
	health := observability.NewHealthChecker()

	ctx := context.Background()

	pgPool, err := repository.NewPostgresPool(ctx, cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	health.Register("postgres", pgPool)
{% if use_redis %}

	redisCache := cache.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	health.Register("redis", redisCache)
	redisLocker := locking.NewRedisLocker(redisCache.Client())
{% endif %}
{% if use_s3 %}

	objectStorage, err := storage.NewS3Storage(ctx, cfg.S3.Endpoint, cfg.S3.Region, cfg.S3.Bucket, cfg.S3.AccessKey, cfg.S3.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("configure s3 storage: %w", err)
	}
{% endif %}
{% if use_kafka %}

	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
{% endif %}

	notifier := notification.NewHTTPNotifier(cfg.Notification.BaseURL, cfg.Notification.APIKey)

	userRepo := repository.NewUserRepository(pgPool)
	userService := service.NewUserService(
		userRepo,
		notifier,
{% if use_redis %}
		redisCache,
{% endif %}
{% if use_kafka %}
		kafkaProducer,
{% endif %}
		logger,
	)
	userHandler := handler.NewUserHandler(userService, dto.NewValidator())

	router := handler.NewRouter(handler.RouterConfig{
		UserHandler:   userHandler,
		HealthChecker: health,
		Metrics:       metrics,
		Logger:        logger,
	})
{% if use_kafka %}

	// The consumer is just another entry point into the same service the
	// HTTP handler calls — here it only logs, but it could just as easily
	// trigger inventory reservation, a notification, etc.
	consumerCtx, cancelConsumer := context.WithCancel(context.Background())
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, events.TopicUserCreated,
		func(ctx context.Context, msg kafkago.Message) error {
			logger.Info("received user.created event", slog.String("key", string(msg.Key)))
			return nil
		},
		logger,
	)
	go kafkaConsumer.Run(consumerCtx)
{% endif %}

	return &Application{
		Config: cfg,
		Logger: logger,
		Router: router,
{% if use_s3 %}
		Storage: objectStorage,
{% endif %}
{% if use_redis %}
		Locker: redisLocker,
{% endif %}
		pgPool: pgPool,
{% if use_redis %}
		redisCache: redisCache,
{% endif %}
{% if use_kafka %}
		kafkaProducer:  kafkaProducer,
		kafkaConsumer:  kafkaConsumer,
		cancelConsumer: cancelConsumer,
{% endif %}
	}, nil
}

// Close releases every resource opened in BuildApplication. Call it with
// defer right after a successful BuildApplication call.
func (a *Application) Close() {
	a.pgPool.Close()
{% if use_redis %}
	_ = a.redisCache.Close()
{% endif %}
{% if use_kafka %}
	a.cancelConsumer()
	_ = a.kafkaProducer.Close()
	_ = a.kafkaConsumer.Close()
{% endif %}
}
