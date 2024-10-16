package http_utils

const (
	SUCCESS               = "SUCCESS"
	CREATED               = "CREATED"
	NOT_FOUND             = "NOT_FOUND"
	INVALID_REQUEST       = "INVALID_REQUEST"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	TIME_OUT              = "TIME_OUT"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewResponse(message string, data interface{}, err interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
		Error:   err,
	}
}

func NewOKResponse(message string, data interface{}) Response {
	return NewResponse(message, data, nil)
}

func NewErrorResponse(message string, err interface{}) Response {
	return NewResponse(message, nil, err)
}

func NewBindingErrorResponse(resourceName, err string) Response {
	return NewErrorResponse(INVALID_REQUEST, err)
}

func NewAttributeErrorResponse(err AttributeError) Response {
	return NewErrorResponse(INVALID_REQUEST, err)
}

type AttributeError struct {
	Attribute  string `json:"attribute"`
	Cause      string `json:"cause"`
	Constraint string `json:"constraint"`
}
