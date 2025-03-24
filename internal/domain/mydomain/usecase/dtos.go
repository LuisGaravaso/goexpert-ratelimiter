package usecase

type UsecaseInputDTO struct {
	Requester string `json:"id"`
}

type UsecaseOutputDTO struct {
	Message string `json:"message"`
}
