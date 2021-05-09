package api

import (
	"findhotel.com/geo-service/service"
	"findhotel.net/geo-api/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/ziflex/lecho/v2"
	"net/http"
)

//StartHttpServer start echo http web server and runs middlewares
func StartHttpServer() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = true
	e.Logger = lecho.From(log.Logger) //Force echo to use zerolog
	//error handler
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	//middlewares
	setHttpMiddlewares(e)
	//set endpoints
	geoService := startGeoService()
	defer geoService.Stop()
	setHttpEndpoints(e, NewGeoLocationApi(geoService))
	//

	httpConfig := &http.Server{
		Addr: fmt.Sprintf(":%s", config.C.HttpPort),
	}
	log.Fatal().Err(e.StartServer(httpConfig)).Msg("")
}

func setHttpMiddlewares(e *echo.Echo) {
	//e.Use(middleware.BodyDump(BodyDumper)) //for debugging
	e.Use(LogMiddleware())
	//config CORS
	if config.C.Cors.Domain == "*" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}))
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{config.C.Cors.Domain},
			MaxAge:       86400,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "cache-control"},
		}))
	}
}

func setHttpEndpoints(e *echo.Echo, geoLocationApi GeoLocationApi) {
	e.GET("/api/v1/geo-location/:ip", geoLocationApi.GetGeoLocationInfo)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

func startGeoService() service.GeoService {
	geoService, err := service.NewGeoService(service.Config{
		MongodbUrl:      config.C.MongoDB.URL,
		MongodbUsername: config.C.MongoDB.Username,
		MongodbPassword: config.C.MongoDB.Password,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to geolocation service")
	}
	return geoService
}
