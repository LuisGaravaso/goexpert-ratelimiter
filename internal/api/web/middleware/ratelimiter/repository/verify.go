package repository

import (
	"ratelim/internal/api/web/middleware/ratelimiter/entity"
	"time"
)

type Store interface {
	SetServiceConfig(entity.ServiceConfig) error
	GetServiceRateLimit(key string) (entity.ServiceConfig, error)
	IncrementRequestCount(key string, windowKey string) (int, error)
	SetExpiration(key string, windowKey string, ttl time.Duration) error
}
