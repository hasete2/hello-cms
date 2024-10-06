package models

type Response struct {
	StatusCode int16  `json:"status_code"`
	Message    string `json:"message"`
}
