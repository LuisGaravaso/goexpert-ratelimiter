package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	godotenv.Load(".env")
	os.Setenv("RATE_LIMIT_CONFIG_PATH", "../../configs/middleware/services.yaml")

	go func() {
		router := NewRouter()
		err := router.Run(":8081")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second) // espera o servidor subir

	code := m.Run()
	os.Exit(code)
}

func TestRateLimiter_SimpleSequential(t *testing.T) {
	client := &http.Client{Timeout: 2 * time.Second}

	headersList := []map[string]string{
		{"Token-Key": ""},        // Default
		{"Api-Key": "abcd1234"},  // Service A
		{"Api-Key": "efgh5678"},  // Bloqueado
		{"Api-Key": "ijkl91011"}, // Alta RPS
		{"Api-Key": "mnop1213"},  // Desconhecido
	}

	expectedStatusFirst := []int{
		http.StatusOK,        // Default deve aceitar a 1ª
		http.StatusOK,        // Service A também
		http.StatusForbidden, // Serviço explicitamente bloqueado
		http.StatusOK,        // Alta RPS
		http.StatusOK,        // Desconhecido, aplica default
	}

	expectedLastStatus := []int{
		http.StatusTooManyRequests, // Default deve bloquear
		http.StatusTooManyRequests, // Service A
		http.StatusForbidden,       // Serviço explicitamente bloqueado
		http.StatusOK,              // Alta RPS
		http.StatusTooManyRequests, // Desconhecido
	}

	totalRequests := 50

	for i, headers := range headersList {

		var firstStatus, lastStatus int

		for j := 1; j <= totalRequests; j++ {
			req, _ := http.NewRequest("GET", "http://localhost:8081/hello", nil)
			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err := client.Do(req)
			assert.NoError(t, err)

			if j == 1 {
				firstStatus = resp.StatusCode
			}
			if j == totalRequests {
				lastStatus = resp.StatusCode
			}

			resp.Body.Close()
		}

		assert.Equal(t, expectedStatusFirst[i], firstStatus, "Serviço %d - Primeiro status", i+1)
		assert.Equal(t, expectedLastStatus[i], lastStatus, "Serviço %d - Último status", i+1)
	}
}
