package verify

import "context"

type VerifyUsecaseInterface interface {
	Verify(ctx context.Context, input VerifyInputDTO) VerifyOutputDTO
}
