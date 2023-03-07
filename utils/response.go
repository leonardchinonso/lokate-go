package utils

import "net/http"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseOK(message string, data interface{}) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}
