package usecase

import "context"

type MydomainUsecaseInterface interface {
	Execute(ctx context.Context, input UsecaseInputDTO) (UsecaseOutputDTO, error)
}
