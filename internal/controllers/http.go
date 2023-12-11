package controllers

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
