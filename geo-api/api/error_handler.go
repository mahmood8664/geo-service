package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

//AppError is api response in case of error
//goland:noinspection GoNameStartsWithPackageName
type AppError struct {
	ErrorCause    error             `json:"-"`
	HttpErrorCode int               `json:"-"`
	ErrorMessage  string            `json:"error_message,omitempty"`
	ErrorData     map[string]string `json:"error_data,omitempty" example:"errorParam:value"`
}

func (e *AppError) Error() string {
	if e.ErrorCause != nil {
		return fmt.Sprintf("AppError with message '%s' and http code '%d' and cause '%s'",
			e.ErrorMessage, e.HttpErrorCode, e.ErrorCause.Error())
	} else {
		return fmt.Sprintf("AppError with message '%s' and http code '%d'",
			e.ErrorMessage, e.HttpErrorCode)
	}
}

//////////////////////////////////////////////////////////
// Define your custom error HERE:

//ParseError parse error base on the error message. In some libraries, in case of error they return error string. we can
//parse them based on the error message they return.
func ParseError(err error) error {
	if err == nil {
		return nil
	}

	var e = &AppError{}
	if errors.As(err, &e) {
		return e
	} else {
		return parseErrorMessage(err)
	}
}

func parseErrorMessage(err error) error {

	switch {
	case err.Error() == "mongo: no documents in result" || err.Error() == "data not found":
		return &AppError{
			ErrorCause:    err,
			HttpErrorCode: http.StatusNotFound,
			ErrorMessage:  "Not Found",
			ErrorData:     nil,
		}
	default:
		log.Error().Err(err).Msg("cannot parse error")
		return &AppError{
			ErrorCause:    err,
			HttpErrorCode: http.StatusInternalServerError,
			ErrorMessage:  "Unknown Error",
			ErrorData:     nil,
		}
	}
}

func InternalServerError(msg string) error {
	return &AppError{
		HttpErrorCode: http.StatusInternalServerError,
		ErrorMessage:  msg,
	}
}

func BadRequest(msg string) error {
	return &AppError{
		HttpErrorCode: http.StatusBadRequest,
		ErrorMessage:  msg,
	}
}

//CustomHTTPErrorHandler handles error based on the type of error. If it is type of AppError, it converts it into json
//but for other kind of error, first will convert it to AppError
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if !c.Response().Committed {
		if e, ok := err.(*AppError); ok {
			handleError(c, e)
		} else if e, ok := err.(*echo.HTTPError); ok {
			handleHttpError(c, e)
		} else {
			err := ParseError(err)
			handleError(c, err.(*AppError))
		}
	}
}

func handleError(c echo.Context, e *AppError) {
	if err2 := c.JSON(e.HttpErrorCode, e); err2 != nil {
		log.Error().Msg("error in converting AppError to json: " + err2.Error())
	}
}

func handleHttpError(c echo.Context, e *echo.HTTPError) {
	response := AppError{
		ErrorMessage: fmt.Sprint(e.Message),
	}

	if err2 := c.JSON(e.Code, response); err2 != nil {
		log.Error().Msg("error in converting HttpError to json: " + err2.Error())
	}
}
