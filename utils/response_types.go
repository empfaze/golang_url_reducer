package utils

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	STATUS_OK    = "OK"
	STATUS_ERROR = "Error"
)

func OK() Response {
	return Response{
		Status: STATUS_OK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: STATUS_ERROR,
		Error:  msg,
	}
}
