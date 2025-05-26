package handlers

type GetQuotesResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type QuitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
