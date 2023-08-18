package entity

type Response struct {
	Code    int     `json:"code"`
	Message any     `json:"message"`
	Data    any     `json:"data"`
	Errors  []error `json:"errors"`
}
