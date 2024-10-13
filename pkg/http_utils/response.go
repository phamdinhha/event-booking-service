package http_utils

import (
	"net/http"
)

const (
	SUCCESS                            = "SUCCESS"
	CREATED                            = "CREATED"
	NOT_FOUND                          = "NOT_FOUND"
	DUPLICATED_SAMPLE                  = "DUPLICATED_SAMPLE"
	DUPLICATED_SUBJECT                 = "DUPLICATED_SUBJECT"
	INVALID_REQUEST                    = "INVALID_REQUEST"
	INTERNAL_SERVER_ERROR              = "INTERNAL_SERVER_ERROR"
	TIME_OUT                           = "TIME_OUT"
	SAMPLE_NOT_FOUND                   = "SAMPLE_NOT_FOUND"
	PRIMARY_SAMPLE_ALREADY_EXISTS      = "PRIMARY_SAMPLE_ALREADY_EXISTS"
	PRIMARY_SAMPLE_FORBIDDEN_TO_DELETE = "PRIMARY_SAMPLE_FORBIDDEN_TO_DELETE"
)

type Response struct {
	Success bool   `json:"success"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type ResponseOK struct {
	Response
	Data interface{} `json:"data"`
}

type StringErrorResponse struct {
	Response
	Error string `json:"error"`
}

type MapErrorResponse struct {
	Response
	Error map[string]string `json:"error"`
}

type AttributeErrorResponse struct {
	Response
	Error AttributeError `json:"error"`
}

type AttributeError struct {
	Attribute  string `json:"attribute"`
	Cause      string `json:"cause"`
	Constraint string `json:"constraint"`
}

func NewBindingErrorResponse(
	resource_name string,
	err string,
) *StringErrorResponse {
	return &StringErrorResponse{
		Response: Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: INVALID_REQUEST,
		},
		Error: err,
	}
}

func NewAttributeErrorResponse(
	err AttributeError,
) *AttributeErrorResponse {
	return &AttributeErrorResponse{
		Response: Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: INVALID_REQUEST,
		},
		Error: err,
	}
}

func NewOKResponse(
	status_code int64,
	message string,
	data interface{},
) *ResponseOK {
	var to_return ResponseOK

	response := Response{
		Success: true,
		Code:    status_code,
		Message: message,
	}
	to_return = ResponseOK{
		Response: response,
		Data:     data,
	}
	return &to_return
}

func NewErrorResponse(
	status_code int64,
	message string,
	err interface{},
) interface{} {
	var to_return interface{}
	response := Response{
		Success: false,
		Code:    status_code,
		Message: message,
	}

	switch err.(type) {
	case string:
		to_return = &StringErrorResponse{
			Response: response,
			Error:    err.(string),
		}
	case map[string]string:
		to_return = &MapErrorResponse{
			Response: response,
			Error:    err.(map[string]string),
		}

	}
	return to_return
}
