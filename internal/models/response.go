package models

type Response struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Headers []Params `json:"headers"`
	Body    string   `json:"body"`
}
