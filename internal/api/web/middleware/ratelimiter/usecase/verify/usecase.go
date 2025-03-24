package verify

import (
	"context"
	"fmt"
	"net/http"
	"ratelim/internal/api/web/middleware/ratelimiter/repository"
	"time"
)

type VerifyUsecase struct {
	RateLimiterRepository repository.Store
}

func NewVerifyUsecase(rateLimiterRepository repository.Store) *VerifyUsecase {
	return &VerifyUsecase{
		RateLimiterRepository: rateLimiterRepository,
	}
}

func (v *VerifyUsecase) Verify(ctx context.Context, input VerifyInputDTO) VerifyOutputDTO {
	// Obter chave de rate limit
	var key string
	if input.ApiKey == "" {
		key = input.ClientIp
	} else {
		key = input.ApiKey
	}

	// Obter config de rate limit do repositório
	config, err := v.RateLimiterRepository.GetServiceRateLimit(key)
	if err != nil {
		return VerifyOutputDTO{
			Key:     key,
			Name:    config.Name,
			Blocked: true,
			Message: "Erro ao buscar configuração de rate limit",
			Status:  http.StatusInternalServerError,
		}
	}

	// Retornar se o serviço estiver explicitamente bloqueado
	if !config.Valid {
		return VerifyOutputDTO{
			Key:     key,
			Name:    config.Name,
			Blocked: true,
			Message: "Serviço bloqueado",
			Status:  http.StatusForbidden,
		}
	}

	// Calcular janela atual
	windowSeconds := config.AllowedRPS
	if windowSeconds <= 0 {
		windowSeconds = 60 // fallback seguro se não estiver configurado corretamente
	}

	now := time.Now()
	windowSize := int64(config.AllowedRPS)
	if windowSize == 0 {
		windowSize = 60 // segurança
	}
	windowTimestamp := now.Unix() / windowSize
	windowKey := fmt.Sprintf("%d", windowTimestamp)

	// Incrementar contador
	count, err := v.RateLimiterRepository.IncrementRequestCount(config.Key, windowKey)
	if err != nil {
		return VerifyOutputDTO{
			Key:     key,
			Name:    config.Name,
			Blocked: true,
			Message: "Erro interno ao contar requisições",
			Status:  http.StatusInternalServerError,
		}
	}

	// Aplicar TTL apenas se for a primeira requisição da janela
	if count == 1 {
		ttl := time.Duration(windowSize+5) * time.Second
		_ = v.RateLimiterRepository.SetExpiration(config.Key, windowKey, ttl)
	}

	// Verificar se está bloqueado
	if count > config.AllowedRPS {
		// Calcular quando o bloqueio termina
		windowResetAt := time.Unix((windowTimestamp+1)*windowSize, 0)
		msg := fmt.Sprintf(
			"Rate limit excedido para o serviço '%s': %d requisições permitidas por segundo. Bloqueado até %s.",
			config.Name,
			config.AllowedRPS,
			windowResetAt.Format("15:04:05"),
		)
		return VerifyOutputDTO{
			Key:     key,
			Name:    config.Name,
			Blocked: true,
			Message: msg,
			Status:  http.StatusTooManyRequests,
		}
	}

	// Requisição liberada
	return VerifyOutputDTO{
		Key:     key,
		Name:    config.Name,
		Blocked: false,
		Message: "",
		Status:  http.StatusOK,
	}
}
