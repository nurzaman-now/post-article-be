package handlers

import "github.com/gin-gonic/gin"

// APIResponse adalah struktur standar response untuk REST API
type APIResponse struct {
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// ResponseSuccess mengirim HTTP 200 OK dengan format response sukses yang seragam
func ResponseSuccess(ctx *gin.Context, data any, message string) {
	ctx.JSON(200, APIResponse{
		Status:  "success",
		Data:    data,
		Message: message,
	})
}

// ResponseError mengirim HTTP 400 Bad Request dengan format response error yang seragam
func ResponseError(ctx *gin.Context, data any, message string) {
	ctx.JSON(400, APIResponse{
		Status:  "error",
		Data:    data,
		Message: message,
	})
}

// ResponseCustom mengirim HTTP status code kustom dengan format APIResponse yang seragam
func ResponseCustom(ctx *gin.Context, httpCode int, status string, data any, message string) {
	ctx.JSON(httpCode, APIResponse{
		Status:  status,
		Data:    data,
		Message: message,
	})
}
