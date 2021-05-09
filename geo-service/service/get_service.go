package service

import (
	"errors"
	config2 "findhotel.com/geo-service/config"
	mongodb "findhotel.com/geo-service/db"
	"findhotel.com/geo-service/model"
	"github.com/rs/zerolog/log"
)

var NotFoundError = errors.New("data not found")

type Config struct {
	MongodbUrl      string
	MongodbUsername string
	MongodbPassword string
}

var geoService = geoServiceImpl{connect: false}

type GeoService interface {
	//GetGeoLocation return model.GeoLocation data of the given IP address, return error if not successful.
	//Just Start GeoService before using this service. If data not found will return NotFoundError
	GetGeoLocation(ip string) (model.GeoLocation, error)
	//Stop disconnect from database and free resources.
	Stop()
}

//NewGeoService try to connect to DB and provide geo services or return error if cannot connect to database
func NewGeoService(config Config) (GeoService, error) {
	if geoService.connect {
		return &geoService, nil
	}
	config2.Init()
	config2.C.MongodbUrl = config.MongodbUrl
	config2.C.MongodbUsername = config.MongodbUsername
	config2.C.MongodbPassword = config.MongodbPassword
	//
	clientWrapper := mongodb.NewClientWrapperImpl()
	err := clientWrapper.Connect()
	if err != nil {
		log.Error().Err(err).Msg("cannot connect to database")
		return &geoService, err
	}
	//
	geoService.dao = mongodb.NewGeoLocationDaoImpl(clientWrapper)
	geoService.connect = true
	geoService.clientWrapper = clientWrapper
	return &geoService, nil
}

type geoServiceImpl struct {
	dao           mongodb.GeoLocationDao
	clientWrapper mongodb.ClientWrapper
	connect       bool
	config        Config
}

func (r *geoServiceImpl) Stop() {
	if r.connect {
		r.clientWrapper.Close()
	}
	r.connect = false
}

func (r *geoServiceImpl) GetGeoLocation(ip string) (model.GeoLocation, error) {
	geoLocation := model.GeoLocation{}
	if !r.connect {
		return geoLocation, errors.New("service is not started")
	}
	geoLocation, err := r.dao.GetOne(ip)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return geoLocation, NotFoundError
	}
	return geoLocation, err
}
