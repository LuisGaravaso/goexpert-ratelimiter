package usecase_test

import (
	"context"
	"testing"

	"ratelim/internal/domain/mydomain/usecase"

	"github.com/stretchr/testify/assert"
)

// TestMydomainUsecase_ImplementInterface tests if MydomainUsecase implements MydomainUsecaseInterface
// This is a good practice to ensure that the interface is implemented correctly
// that way you can refactor the usecase without breaking the interface
func TestMydomainUsecase_ImplementInterface(t *testing.T) {
	var _ usecase.MydomainUsecaseInterface = &usecase.MydomainUsecase{}
}

// Now Testing the Execute method of MydomainUsecase
func TestMydomainUsecase_Execute_MustPass(t *testing.T) {

	// Arrange
	u := usecase.NewMydomainUsecase()
	input := usecase.UsecaseInputDTO{
		Requester: "World",
	}

	// Act
	output, err := u.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World", output.Message)
}

func TestMydomainUsecase_Execute_MustFailForEmptyRequester(t *testing.T) {

	// Arrange
	u := usecase.NewMydomainUsecase()
	input := usecase.UsecaseInputDTO{
		Requester: "",
	}

	// Act
	output, err := u.Execute(context.Background(), input)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "requester is empty", err.Error())
	assert.Equal(t, "", output.Message)
}
