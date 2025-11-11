package dto

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type DataResponse struct {
	Data interface{} `json:"data,omitempty"`
}
