package handlers

import (
	"net/http"
	"ratelim/internal/domain/mydomain/usecase"

	"github.com/gin-gonic/gin"
)

type HelloService struct {
	usecase usecase.MydomainUsecaseInterface
}

func NewHelloService(usecase usecase.MydomainUsecaseInterface) *HelloService {
	return &HelloService{usecase: usecase}
}

func (h *HelloService) Hello(c *gin.Context) {
	input := usecase.UsecaseInputDTO{
		Requester: c.GetString("Requester"),
	}
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}
