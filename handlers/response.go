package handlers

import "github.com/gin-gonic/gin"

// APIResponse adalah struktur standar response untuk REST API
type APIResponse struct {
	Status  string `json:"status"`
	Meta    any    `json:"meta"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// ResponseSuccess mengirim HTTP 200 OK dengan format response sukses yang seragam
func ResponseSuccess(ctx *gin.Context, meta any, data any, message string, code ...int) {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	ctx.JSON(statusCode, APIResponse{
		Status:  "success",
		Data:    data,
		Message: message,
		Meta:    meta,
	})
}

// ResponseError mengirim HTTP 400 Bad Request dengan format response error yang seragam
func ResponseError(ctx *gin.Context, data any, message string, code ...int) {
	statusCode := 400
	if len(code) > 0 {
		statusCode = code[0]
	}
	ctx.JSON(statusCode, APIResponse{
		Status:  "error",
		Data:    data,
		Message: message,
	})
}

// ResponseCustom mengirim HTTP status code kustom dengan format APIResponse yang seragam
func ResponseValidation(ctx *gin.Context, err error) {
	reqErr := GetErrorMap(err)
	reqMessages := MapToString(reqErr)
	ctx.JSON(422, APIResponse{
		Status:  "error",
		Data:    nil,
		Message: "Data yang di inputkan tidak sesuai: " + reqMessages,
	})
}
