module {{ module_name }}

go {{ go_version }}

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-playground/validator/v10 v10.22.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/prometheus/client_golang v1.20.5
{% if use_redis %}
	github.com/redis/go-redis/v9 v9.6.1
{% endif %}
{% if use_kafka %}
	github.com/segmentio/kafka-go v0.4.47
{% endif %}
{% if use_s3 %}
	github.com/aws/aws-sdk-go-v2 v1.32.2
	github.com/aws/aws-sdk-go-v2/config v1.28.1
	github.com/aws/aws-sdk-go-v2/credentials v1.17.42
	github.com/aws/aws-sdk-go-v2/service/s3 v1.66.1
{% endif %}
)
