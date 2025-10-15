package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Константы типов ответа
const (
	Object = "object" // Используется когда ответ содержит один объект
	Array  = "array"  // Используется когда ответ содержит массив объектов
	Empty  = "empty"  // Используется когда ответ не содержит данных
)

// Response - интерфейс для стандартных ответов API
type Response interface {
	ErrorResponse(c *gin.Context, err error, statusCode int, message string, showError bool)
	ResultResponse(c *gin.Context, message string, dataType string, data interface{})
}

// ErrorResponse - возвращает стандартизированный ответ с ошибкой
func (h *Handler) ErrorResponse(c *gin.Context, err error, statusCode int, message string, showError bool) {
	errorMessage := message
	if showError && err != nil {
		errorMessage = message + ": " + err.Error()
	}

	c.JSON(statusCode, gin.H{
		"status": "error",
		"error": gin.H{
			"code":    statusCode,
			"message": errorMessage,
		},
	})
}

// ResultResponse - возвращает JSON с данными
func (h *Handler) ResultResponse(c *gin.Context, message string, dataType string, data interface{}) {
	response := gin.H{
		"status":  "success",
		"message": message,
		"type":    dataType,
	}

	if data != nil {
		response["data"] = data
	}

	c.JSON(http.StatusOK, response)
}

// BadRequest - возвращает ошибку 400
func (h *Handler) BadRequest(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
}

// InternalError - возвращает ошибку 500
func (h *Handler) InternalError(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusInternalServerError, errors.InternalServerError, false)
}

// NotFound - возвращает ошибку 404
func (h *Handler) NotFound(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusNotFound, errors.NotFound, true)
}

// RespondWithError - вспомогательная функция для возврата ошибки (для чистого http)
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON - вспомогательная функция для возврата JSON (для чистого http)
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ResultResponse represents a successful API response
// @Description Standard wrapper for successful responses
type ResultResponse struct {
	// Overall status of the response
	// @example "success"
	Status string `json:"status"`
	// Response payload container
	Response struct {
		// Human-readable success message
		// @example "Patient created successfully"
		Message string `json:"message"`
		// Type of data: "object", "array", or "empty"
		// @example "object"
		// @enum object array empty
		Type string `json:"type"`
		// Actual data (optional)
		Data interface{} `json:"data,omitempty"`
	} `json:"response"`
}

// ResultError represents a standard API error response
// @Description Generic error response format used across the API
type ResultError struct {
	// Always "error" for error responses
	// @example "error"
	Status string `json:"status"`
	// Error details
	Response struct {
		// HTTP status code
		// @example 400
		Code int `json:"code"`
		// Human-readable error message
		// @example "Invalid request body"
		Message string `json:"message"`
	} `json:"response"`
}

// InternalServerError
type InternalServerError struct {
	Status   string `json:"status" example:"InternalServerError"` // error
	Response struct {
		Code    int    `json:"code" example:"500"`
		Message string `json:"message" example:"Внутренняя ошибка сервера"`
	} `json:"response"`
}

// IncorrectFormatError
type IncorrectFormatError struct {
	Status   string `json:"status" example:"IncorrectFormatError"` // error
	Response struct {
		Code    int    `json:"code" example:"400"`
		Message string `json:"message" example:"Неверный формат запроса"`
	} `json:"response"`
}

// IncorrectDataError
type IncorrectDataError struct {
	Status   string `json:"status" example:"IncorrectDataError"` // error
	Response struct {
		Code    int    `json:"code" example:"401"`
		Message string `json:"message" example:"Некорректные данные"`
	} `json:"response"`
}

// NotFoundError
type NotFoundError struct {
	Status   string `json:"status" example:"NotFoundError"` // error
	Response struct {
		Code    int    `json:"code" example:"404"`
		Message string `json:"message" example:"Данные не найдены"`
	} `json:"response"`
}

// ValidationError
type ValidationError struct {
	Status   string `json:"status" example:"ValidationError"` // error
	Response struct {
		Code    int    `json:"code" example:"422"`
		Message string `json:"message" example:"Ошибка валидации"`
	} `json:"response"`
}
