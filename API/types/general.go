package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessfulResponse struct {
	Message string `json:"message"`
}
