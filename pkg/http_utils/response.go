package http_utils

import (
	"net/http"
)

const (
	SUCCESS               = "SUCCESS"
	CREATED               = "CREATED"
	NOT_FOUND             = "NOT_FOUND"
	DUPLICATED_SAMPLE     = "DUPLICATED_SAMPLE"
	DUPLICATED_SUBJECT    = "DUPLICATED_SUBJECT"
	INVALID_REQUEST       = "INVALID_REQUEST"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	TIME_OUT              = "TIME_OUT"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewResponse(statusCode int, message string, data interface{}, err interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
		Error:   err,
	}
}

func NewOKResponse(statusCode int, message string, data interface{}) Response {
	return NewResponse(statusCode, message, data, nil)
}

func NewErrorResponse(statusCode int, message string, err interface{}) Response {
	return NewResponse(statusCode, message, nil, err)
}

func NewBindingErrorResponse(resourceName, err string) Response {
	return NewErrorResponse(http.StatusBadRequest, INVALID_REQUEST, err)
}

func NewAttributeErrorResponse(err AttributeError) Response {
	return NewErrorResponse(http.StatusBadRequest, INVALID_REQUEST, err)
}

type AttributeError struct {
	Attribute  string `json:"attribute"`
	Cause      string `json:"cause"`
	Constraint string `json:"constraint"`
}
