package api

import (
	"findhotel.com/geo-service/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GeoLocationApi interface {
	GetGeoLocationInfo(ctx echo.Context) error
}

type GeoLocationApiImpl struct {
	geoService service.GeoService
}

func NewGeoLocationApi(geoService service.GeoService) GeoLocationApiImpl {
	return GeoLocationApiImpl{
		geoService: geoService,
	}
}

// GetGeoLocationInfo Get GeoLocation Data
// @Summary Get Geolocation Data
// @Description This API try to find geolocation of the given IP address
// @Tags Geo
// @Produce json
// @Param ip path string true " "
// @Success 200 {object} GeoLocationResponse "GeoLocationResponse"
// @Failure 400,404,500 {object} AppError
// @Router /api/v1/geo-location/{ip} [get]
func (r GeoLocationApiImpl) GetGeoLocationInfo(ctx echo.Context) error {
	ip := ctx.Param("ip")
	if ip == "" {
		log.Warn().Msg("Bad request")
		return BadRequest("ip address must has value")
	}
	geoLocation, err := r.geoService.GetGeoLocation(ip)
	if err != nil {
		return ParseError(err)
	}
	return ctx.JSON(http.StatusOK, FromGeoLocationModel(geoLocation))
}
