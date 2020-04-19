package errors

import (
	"net/http"
)

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func InternalServerError(msg string) ErrorResponse {
	if msg == "" {
		msg = "Error while processing your request."
	}
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

func NotFound(msg string) ErrorResponse {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

func Unauthorized(msg string) ErrorResponse {
	if msg == "" {
		msg = "Unauthorized request."
	}
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}

func Forbidden(msg string) ErrorResponse {
	if msg == "" {
		msg = "Forbidden request."
	}
	return ErrorResponse{
		Status:  http.StatusForbidden,
		Message: msg,
	}
}

func BadRequest(msg string) ErrorResponse {
	if msg == "" {
		msg = "Bad request."
	}
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}
