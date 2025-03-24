package verify

type VerifyInputDTO struct {
	ApiKey   string `json:"api_key"`
	ClientIp string `json:"client_ip"`
}

type VerifyOutputDTO struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Blocked bool   `json:"blocked"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
