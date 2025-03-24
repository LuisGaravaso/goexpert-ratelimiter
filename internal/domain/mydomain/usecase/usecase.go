package usecase

import (
	"context"
	"errors"
)

type MydomainUsecase struct{}

func NewMydomainUsecase() *MydomainUsecase {
	return &MydomainUsecase{}
}

func (u *MydomainUsecase) Execute(ctx context.Context, input UsecaseInputDTO) (UsecaseOutputDTO, error) {

	r := input.Requester
	if r == "" {
		return UsecaseOutputDTO{}, errors.New("requester is empty")
	}

	return UsecaseOutputDTO{
		Message: "Hello, " + r,
	}, nil
}
