package api

import (
	"encoding/json"
	"findhotel.com/geo-service/model"
	"findhotel.com/geo-service/service"
	"findhotel.net/geo-api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type GeoServiceMock struct {
	mock.Mock
}

func (r *GeoServiceMock) GetGeoLocation(ip string) (model.GeoLocation, error) {
	args := r.Called(ip)
	return args.Get(0).(model.GeoLocation), args.Error(1)
}

func (r *GeoServiceMock) Stop() {
	r.Called()
}

func TestGeoLocationApi_GetGeoLocationInfo_OK(t *testing.T) {

	geoServiceMock := GeoServiceMock{}
	location1 := model.GeoLocation{
		IP:           "192.1687.1.1",
		CountryCode:  "AA",
		Country:      "ABC",
		Longitude:    10,
		Latitude:     20,
		City:         "CITY",
		MysteryValue: "M",
		Id:           nil,
	}
	geoServiceMock.On("GetGeoLocation", "192.168.1.1").Return(location1, nil)

	config.Init("../resources/config.yml")

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/geo-location/:ip", &strings.Reader{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	//
	c.SetParamNames("ip")
	c.SetParamValues("192.168.1.1")
	//
	h := NewGeoLocationApi(&geoServiceMock)
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	setHttpMiddlewares(e)
	setHttpEndpoints(e, h)

	// Assertions
	if assert.NoError(t, h.GetGeoLocationInfo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		marshal, err := json.Marshal(FromGeoLocationModel(location1))
		assert.NoError(t, err)

		assert.Equal(t, string(marshal)+"\n", rec.Body.String())
	}

}

func TestGeoLocationApi_GetGeoLocationInfo_NotFound(t *testing.T) {

	geoServiceMock := GeoServiceMock{}
	geoServiceMock.On("GetGeoLocation", "192.168.1.2").Return(model.GeoLocation{}, service.NotFoundError)

	config.Init("../resources/config.yml")

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/geo-location/:ip", &strings.Reader{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	//
	c.SetParamNames("ip")
	c.SetParamValues("192.168.1.2")
	//
	h := NewGeoLocationApi(&geoServiceMock)
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	setHttpMiddlewares(e)
	setHttpEndpoints(e, h)

	// Assertions
	err := h.GetGeoLocationInfo(c)
	if assert.Error(t, err) {
		e, ok := err.(*AppError)
		assert.Equal(t, ok, true)
		assert.Equal(t, e.ErrorMessage, "Not Found")
		assert.Equal(t, e.HttpErrorCode, 404)
	}

}

func TestGeoLocationApi_GetGeoLocationInfo_BadRequest(t *testing.T) {

	geoServiceMock := GeoServiceMock{}
	geoServiceMock.On("GetGeoLocation", "192.168.1.2").Return(model.GeoLocation{}, service.NotFoundError)

	config.Init("../resources/config.yml")

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/geo-location/", &strings.Reader{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	//
	h := NewGeoLocationApi(&geoServiceMock)
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	setHttpMiddlewares(e)
	setHttpEndpoints(e, h)

	// Assertions
	err := h.GetGeoLocationInfo(c)
	if assert.Error(t, err) {
		e, ok := err.(*AppError)
		assert.Equal(t, ok, true)
		assert.Equal(t, e.ErrorMessage, "ip address must has value")
		assert.Equal(t, e.HttpErrorCode, 400)
	}

}
