package service

import (
	"errors"
	"findhotel.com/geo-service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type GeoLocationDaoMock struct {
	mock.Mock
}

func (r *GeoLocationDaoMock) GetOne(ip string) (geoLocation model.GeoLocation, err error) {
	args := r.Called(ip)
	return args.Get(0).(model.GeoLocation), args.Error(1)
}

func (r *GeoLocationDaoMock) Insert(geoLocation model.GeoLocation) (id string, err error) {
	args := r.Called(geoLocation)
	return args.String(0), args.Error(1)
}

func (r *GeoLocationDaoMock) InsertMany(geoLocation []model.GeoLocation) (ids []interface{}, err error) {
	args := r.Called(geoLocation)
	return args.Get(0).([]interface{}), args.Error(1)
}

////

type ClientWrapperMock struct {
	mock.Mock
}

func (r *ClientWrapperMock) Client() *mongo.Client {
	args := r.Called()
	return args.Get(0).(*mongo.Client)
}

func (r *ClientWrapperMock) Connect() error {
	r.Called()
	return nil
}

func (r *ClientWrapperMock) Close() {
	r.Called()
}

/////

func TestGeoServiceImpl_GetOne(t *testing.T) {

	geoDaoMock := GeoLocationDaoMock{}

	geoDaoMock.On("GetOne", "192.168.1.1").
		Return(model.GeoLocation{Id: "1", IP: "192.168.1.1", CountryCode: "IR", Country: "Iran", Longitude: 20, City: "Teh", Latitude: 10, MysteryValue: "123"}, nil)
	geoDaoMock.On("GetOne", "192.168.1.2").Return(model.GeoLocation{}, errors.New("mongo: no documents in result"))

	///
	clientWrapperMock := ClientWrapperMock{}
	clientWrapperMock.On("Connect").Return()
	clientWrapperMock.On("Close").Return()
	//
	geoService = geoServiceImpl{ //We cannot use NewGeoService() because it tries to connect to real database, so we just construct service ourselves
		connect:       true,
		clientWrapper: &clientWrapperMock,
		dao:           &geoDaoMock,
		config: Config{
			MongodbUrl:      "aaaa",
			MongodbUsername: "uuuu",
			MongodbPassword: "ppp",
		},
	}

	one, _ := geoService.GetGeoLocation("192.168.1.1")
	_, err2 := geoService.GetGeoLocation("192.168.1.2")

	assert.Equal(t, one, model.GeoLocation{Id: "1", IP: "192.168.1.1", CountryCode: "IR", Country: "Iran", Longitude: 20, City: "Teh", Latitude: 10, MysteryValue: "123"})
	assert.Error(t, err2)
	assert.Equal(t, err2.Error(), "data not found")
	geoDaoMock.MethodCalled("GetOne", "192.168.1.1")
	geoDaoMock.MethodCalled("GetOne", "192.168.1.2")
}
