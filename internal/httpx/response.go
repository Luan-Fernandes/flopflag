package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func OK(c *gin.Context, status int, message string, data any) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, message string, errs any) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func OK200(c *gin.Context, message string, data any) {
	OK(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data any) {
	OK(c, http.StatusCreated, message, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func BadRequest(c *gin.Context, message string, errs any) {
	Fail(c, http.StatusBadRequest, message, errs)
}

func Unauthorized(c *gin.Context, message string, errs any) {
	Fail(c, http.StatusUnauthorized, message, errs)
}

func Forbidden(c *gin.Context, message string, errs any) {
	Fail(c, http.StatusForbidden, message, errs)
}

func NotFound(c *gin.Context, message string, errs any) {
	Fail(c, http.StatusNotFound, message, errs)
}

func Conflict(c *gin.Context, message string, errs any) {
	Fail(c, http.StatusConflict, message, errs)
}

func InternalError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, message, nil)
}

type AppError struct {
	Status  int
	Message string
	Errors  any
}

func (e *AppError) Error() string { return e.Message }

func NewAppError(status int, message string, errs any) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
		Errors:  errs,
	}
}

func FailErr(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		Fail(c, appErr.Status, appErr.Message, appErr.Errors)
		return
	}

	InternalError(c, "erro interno")
}
