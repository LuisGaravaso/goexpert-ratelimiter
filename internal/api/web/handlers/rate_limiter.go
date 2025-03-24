package handlers

import (
	v "ratelim/internal/api/web/middleware/ratelimiter/usecase/verify"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	usecase v.VerifyUsecaseInterface
}

func NewRateLimiter(usecase v.VerifyUsecaseInterface) *RateLimiter {
	return &RateLimiter{usecase: usecase}
}

func (r *RateLimiter) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.GetHeader("Api-Key")
		client_ip := c.ClientIP()

		input := v.VerifyInputDTO{
			ApiKey:   api_key,
			ClientIp: client_ip,
		}
		block := r.usecase.Verify(c.Request.Context(), input)
		if block.Blocked {
			c.JSON(block.Status, gin.H{"message": block.Message})
			c.AbortWithStatus(block.Status)
			return
		}

		c.Set("Requester", block.Name)

		c.Next()
	}
}
