package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response interface
type Response interface {
	Status() int
	Err() error

	// header getter
	Header() http.Header
	// header setter
	WithHeader(k, v string) Response
}

// ResponseMsg holds messages and needed data
type ResponseMsg struct {
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// Handler interface
type Handler func(ctx *gin.Context) (interface{}, Response)

// WrapFunc is a helper wrapper to make implementing handlers easier
func WrapFunc(a Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		object, result := a(ctx)

		ctx.Header("Content-Type", "application/json")

		var status int
		if result == nil {
			status = http.StatusOK
		} else {
			h := result.Header()
			for k := range h {
				for _, v := range h.Values(k) {
					ctx.Header(k, v)
				}
			}

			if err := result.Err(); err != nil {
				object = struct {
					Error string `json:"err"`
				}{
					Error: err.Error(),
				}
			}
			status = result.Status()
		}

		ctx.IndentedJSON(status, object)
	}
}

type genericResponse struct {
	status int
	err    error
	header http.Header
}

// Status return response status
func (r genericResponse) Status() int {
	return r.status
}

// Err return response error
func (r genericResponse) Err() error {
	return r.err
}

// Header return response header
func (r genericResponse) Header() http.Header {
	if r.header == nil {
		r.header = http.Header{}
	}
	return r.header
}

// WithHeader add header to response
func (r genericResponse) WithHeader(k, v string) Response {
	if r.header == nil {
		r.header = http.Header{}
	}

	r.header.Add(k, v)
	return r
}

// Created return a created response
func Created() Response {
	return genericResponse{status: http.StatusCreated}
}

// Ok return a ok response
func Ok() Response {
	return genericResponse{status: http.StatusOK}
}

// Error generic error response
func Error(err error, code ...int) Response {
	status := http.StatusInternalServerError
	if len(code) > 0 {
		status = code[0]
	}

	if err == nil {
		err = fmt.Errorf("no message")
	}

	return genericResponse{status: status, err: err}
}

// BadRequest result
func BadRequest(err error) Response {
	return Error(err, http.StatusBadRequest)
}

// InternalServerError result
func InternalServerError(err error) Response {
	return Error(err, http.StatusInternalServerError)
}

func Conflict(err error) Response {
	return Error(err, 409)
}

// NotFound response
func NotFound(err error) Response {
	return Error(err, http.StatusNotFound)
}

// UnAuthorized response
func UnAuthorized(err error) Response {
	return Error(err, http.StatusUnauthorized)
}

// Forbidden response
func Forbidden(err error) Response {
	return Error(err, http.StatusForbidden)
}

// Accepted response
func Accepted() Response {
	return genericResponse{status: http.StatusAccepted}
}
