// API 响应结构体定义
package api

import (
	"net/http"
	"time"
)

// APIResponse 统一 API 响应结构
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	PageNum    int `json:"page_num"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	StartIndex int `json:"start_index"`
	EndIndex   int `json:"end_index"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code:      http.StatusOK,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string, err error) APIResponse {
	response := APIResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	if err != nil {
		response.Error = err.Error()
	}
	return response
}

// BadRequestResponse 400 错误响应
func BadRequestResponse(message string, err error) APIResponse {
	return ErrorResponse(http.StatusBadRequest, message, err)
}

// InternalErrorResponse 500 错误响应
func InternalErrorResponse(message string, err error) APIResponse {
	return ErrorResponse(http.StatusInternalServerError, message, err)
}
