package csv

import (
	"findhotel.com/geo-service/config"
	"findhotel.com/geo-service/model"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestDispatcherImpl_Dispatch(t *testing.T) {

	config.C = config.Config{
		NumberOfWorkers: 1,
		InsertBulkSize:  1,
	}

	geoDaoMock := GeoLocationDaoMock{}

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.10.3", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.2", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.4", CountryCode: "AR", Country: "I", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.5", CountryCode: "AR", Country: "AAAAA BBBBB", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.7", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.9", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 90, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.10", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: -90, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.12", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 10, Longitude: 25, MysteryValue: "'aaa'",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.13", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 10, Longitude: 25, MysteryValue: "'aaaa",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.13", CountryCode: "AR", Country: "Irannn", City: "Tehrannn", Latitude: 10, Longitude: 25, MysteryValue: "'aaaa",
	}}).Return([]interface{}{}, mongo.BulkWriteException{WriteErrors: []mongo.BulkWriteError{
		{
			WriteError: mongo.WriteError{
				Index:   0,
				Code:    1011,
				Message: "Duplicate",
			},
		},
	}})

	dispatcher := NewDispatcherImpl(&geoDaoMock)

	recordChan := make(chan []string)
	finishedCh := make(chan bool)
	dispatcher.Dispatch(recordChan, finishedCh)

	data := MockData()
	for _, record := range data {
		recordChan <- record
	}
	finishedCh <- true
	//wait until workers finish their works
	dispatcher.Wait()
	dispatcher.Print()

	geoDaoMock.AssertNumberOfCalls(t, "InsertMany", 10)
}

func MockData() [][]string {
	return [][]string{
		{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"},
		{"235", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"235.325.36.365", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.10.003", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.1", "A", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.2", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.3", "AR", "  ", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.4", "AR", "I ", "Tehran", "12", "20.236", "AAAA"},               //OK
		{"192.168.1.5", "AR", "  AAAAA BBBBB  ", " Tehran", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.6", "AR", "Iran", "  ", "12", "20.236", "AAAA"},
		{"192.168.1.7", "AR", "Iran", "Tehran  ", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.8", "AR", "Iran", "Tehran", "-91", "20.236", "AAAA"},
		{"192.168.1.9", "AR", "Iran", "Tehran", "90", "20.236", "AAAA"},   //OK
		{"192.168.1.10", "AR", "Iran", "Tehran", "-90", "20.236", "AAAA"}, //OK
		{"192.168.1.11", "AR", "Iran", "Tehran", "10", "25", "  "},
		{"192.168.1.12", "AR", "Iran", "Tehran", "10", "25", "'aaa'"},     //OK
		{"192.168.1.13", "AR", "Iran", "Tehran", "10", "25", "'aaaa"},     //OK
		{"192.168.1.13", "AR", "Irannn", "Tehrannn", "10", "25", "'aaaa"}, //Duplicate
		{"", "", "", "", "", "", ""},
	}
}
