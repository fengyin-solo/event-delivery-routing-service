// Package httpx 提供 HTTP 响应与请求解析的通用工具。
package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Response 统一 API 响应结构。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON 以统一结构写出 JSON 响应。
func JSON(w http.ResponseWriter, status, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Code: code, Message: message, Data: data})
}

// OK 写出 200 成功响应。
func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, 0, "ok", data)
}

// Created 写出 201 创建成功响应。
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, 0, "ok", data)
}

// NoContent 写出 204 无内容响应。
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error 写出统一错误响应。
func Error(w http.ResponseWriter, status, code int, message string) {
	JSON(w, status, code, message, nil)
}

// BadRequest 写出 400 响应。
func BadRequest(w http.ResponseWriter, message string) { Error(w, http.StatusBadRequest, 400, message) }

// Unauthorized 写出 401 响应。
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, 401, message)
}

// Forbidden 写出 403 响应。
func Forbidden(w http.ResponseWriter, message string) { Error(w, http.StatusForbidden, 403, message) }

// NotFound 写出 404 响应。
func NotFound(w http.ResponseWriter, message string) { Error(w, http.StatusNotFound, 404, message) }

// Conflict 写出 409 响应。
func Conflict(w http.ResponseWriter, message string) { Error(w, http.StatusConflict, 409, message) }

// InternalError 写出 500 响应。
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, 500, message)
}

// maxBodySize 限制请求体不超过 1MB，避免恶意大请求体。
const maxBodySize = 1 << 20

// Decode 解析 JSON 请求体，限制 1MB，且只允许单个 JSON 值。
// 多个 JSON 值拼接（如 `{"a":1}{"b":2}`）会被拒绝。
func Decode(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, maxBodySize))
	if err := dec.Decode(dst); err != nil {
		return err
	}
	// 拒绝首个值之后的额外内容，防止拼接 JSON 注入多余记录。
	if dec.More() {
		return errors.New("请求体包含多个 JSON 值")
	}
	return nil
}

// Pagination 分页信息。
type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// PageResult 分页查询结果。
type PageResult struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// PageParams 分页入参。
type PageParams struct {
	Page int
	Size int
}

// ParsePagination 从查询参数解析分页信息，并做默认值与上限约束。
func ParsePagination(r *http.Request, defaultSize, maxSize int) PageParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}
	return PageParams{Page: page, Size: size}
}
