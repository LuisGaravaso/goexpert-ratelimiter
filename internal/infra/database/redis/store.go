package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"ratelim/internal/api/web/middleware/ratelimiter/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStore{client: rdb}
}

func (r *RedisStore) SetServiceConfig(cfg entity.ServiceConfig) error {
	ctx := context.Background()

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// Armazena no Redis Hash "rate_limit_config" com campo = cfg.Key
	return r.client.HSet(ctx, "rate_limit_config", cfg.Key, data).Err()
}

func (r *RedisStore) GetServiceRateLimit(key string) (entity.ServiceConfig, error) {
	ctx := context.Background()

	var cfg entity.ServiceConfig

	// 1. Tenta buscar config da chave normalmente
	val, err := r.client.HGet(ctx, "rate_limit_config", key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(val), &cfg)
		return cfg, err
	}

	// 2. Se não encontrou a chave, aplica comportamento com base no default
	if err == redis.Nil {
		// 2.1 Busca config do default
		defaultVal, derr := r.client.HGet(ctx, "rate_limit_config", "default").Result()
		if derr != nil {
			return cfg, fmt.Errorf("configuração default não encontrada: %v", derr)
		}

		var defaultCfg entity.ServiceConfig
		if err := json.Unmarshal([]byte(defaultVal), &defaultCfg); err != nil {
			return cfg, fmt.Errorf("erro ao deserializar config default: %v", err)
		}

		// 2.2 Aplica a config default para esta nova chave (copiando o conteúdo)
		newCfg := defaultCfg
		newCfg.Key = key
		newCfg.Name = fmt.Sprintf("service-%s", RandomString(12))
		newCfg.Valid = true

		// 2.3 Salva essa nova configuração com base na default
		if setErr := r.SetServiceConfig(newCfg); setErr != nil {
			// Mesmo que falhe ao salvar, ainda retornamos a config aplicada
			fmt.Printf("aviso: falha ao salvar config para nova chave %s: %v\n", key, setErr)
		}

		return newCfg, nil
	}

	// 3. Outros erros reais
	return cfg, err
}

func (r *RedisStore) IncrementRequestCount(key string, windowKey string) (int, error) {
	ctx := context.Background()
	fullKey := fmt.Sprintf("rate_limit_counter:%s:%s", key, windowKey)

	count, err := r.client.Incr(ctx, fullKey).Result()
	return int(count), err
}

func (r *RedisStore) SetExpiration(key string, windowKey string, ttl time.Duration) error {
	ctx := context.Background()
	fullKey := fmt.Sprintf("rate_limit_counter:%s:%s", key, windowKey)

	return r.client.Expire(ctx, fullKey, ttl).Err()
}
