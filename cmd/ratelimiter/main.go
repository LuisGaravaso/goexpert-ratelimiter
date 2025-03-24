package main

import (
	"fmt"
	"os"
	"ratelim/internal/api/web/handlers"
	"ratelim/internal/api/web/middleware/ratelimiter/configs"
	"ratelim/internal/api/web/middleware/ratelimiter/usecase/verify"
	"ratelim/internal/domain/mydomain/usecase"
	"ratelim/internal/infra/database/redis"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := NewRouter()
	router.Run(":8080")
}

// NewRouter returns a configured gin.Engine
func NewRouter() *gin.Engine {
	// Load env
	_ = godotenv.Load("cmd/ratelimiter/.env")

	// Load config
	rateLimConfiPath := os.Getenv("RATE_LIMIT_CONFIG_PATH")
	fmt.Printf("Loading config from: %s\n", rateLimConfiPath)
	config, err := configs.LoadConfig(rateLimConfiPath)
	if err != nil || len(config.Services) == 0 {
		panic("Failed to load services config")
	}

	// Setup Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisRepo := redis.NewRedisStore(redisAddr, redisPassword, redisDB)

	for _, service := range config.Services {
		redisRepo.SetServiceConfig(*service)
	}

	// Setup handlers
	usecase := usecase.NewMydomainUsecase()
	helloService := handlers.NewHelloService(usecase)

	ratelimiterUseCase := verify.NewVerifyUsecase(redisRepo)
	rateLimiter := handlers.NewRateLimiter(ratelimiterUseCase)

	// Build router
	router := gin.New()
	router.Use(rateLimiter.Verify())
	router.GET("/hello", helloService.Hello)

	return router
}
