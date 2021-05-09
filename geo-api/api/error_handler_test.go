package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseError(t *testing.T) {
	appError := AppError{
		HttpErrorCode: 100,
		ErrorMessage:  "error one",
	}
	err := ParseError(&appError)
	assert.Equal(t, err, &appError)

	errStr := errors.New("string error")
	err = ParseError(errStr)
	assert.Equal(t, err, &AppError{
		ErrorCause:    errStr,
		HttpErrorCode: http.StatusInternalServerError,
		ErrorMessage:  "Unknown Error",
		ErrorData:     nil,
	})
}

func TestInternalServerError(t *testing.T) {
	err := InternalServerError("this is internal server error")
	assert.Equal(t, err, &AppError{
		HttpErrorCode: http.StatusInternalServerError,
		ErrorMessage:  "this is internal server error",
	})
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("this is bad request error")
	assert.Equal(t, err, &AppError{
		HttpErrorCode: http.StatusBadRequest,
		ErrorMessage:  "this is bad request error",
	})
}

func TestCustomHTTPErrorHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/geo-location/", &strings.Reader{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	CustomHTTPErrorHandler(&AppError{HttpErrorCode: 500}, c)
	assert.True(t, c.Response().Committed)
	assert.Equal(t, c.Response().Status, 500)
}
