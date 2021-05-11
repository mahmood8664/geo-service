package csv

import (
	"findhotel.com/geo-service/config"
	"findhotel.com/geo-service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

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

func TestImporterImpl_Import(t *testing.T) {
	config.C.FileAddress = "/tmp/csv_tmp.csv"
	config.C.NumberOfWorkers = 1 //set to 1 for easy tracking
	config.C.InsertBulkSize = 1
	///////Create file and put csv data
	create, err := os.Create("/tmp/csv_tmp.csv")
	assert.Equal(t, err, nil)
	//put some csv record
	_, err = create.WriteString("\"ip_address\",\"country_code\",\"country\",\"city\",\"latitude\",\"longitude\",\"mystery_value\"\n")
	assert.Equal(t, err, nil)
	_, err = create.WriteString("\"255.32.63\",\"AR\",\"Iran\",\"Tehran\",\"12\",\"20.236\",\"AAAA\"\n")
	assert.Equal(t, err, nil)
	_, err = create.WriteString("\"192.168.1.3\",\"AR\",\"Iran\",\"Tehran\",\"12\",\"20.236\",\"AAAA\"\n")
	assert.Equal(t, err, nil)
	_, err = create.WriteString("\"192.168.1.4\",\"AR\",\"Iran\",\"Tehran\",\"12\",\"20.236\",\"AAAA\"\n")
	assert.Equal(t, err, nil)
	_, err = create.WriteString("\"192.168.1.5\",\"AR\",\"Iran\",\"Tehran\",\"12\",\"20.236\",\"AAAA\"\n")
	assert.Equal(t, err, nil)
	err = create.Close()
	assert.Equal(t, err, nil)
	///////Close csv

	//////Create and config GeoLocationDaoMock
	geoDaoMock := GeoLocationDaoMock{}

	geoDaoMock.On("CreateIndex").Return(nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.3", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.4", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)

	geoDaoMock.On("InsertMany", []model.GeoLocation{{
		IP: "192.168.1.5", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}}).Return([]interface{}{primitive.NewObjectID()}, nil)
	/////Create and config Mongo Client Wrapper
	clientMock := ClientWrapperMock{}
	clientMock.On("Connect").Return()
	clientMock.On("Close").Return()
	/////
	dispatcher := NewDispatcherImpl(&geoDaoMock)
	importer := NewImporterImpl(&clientMock, dispatcher)

	_ = importer.Import()
	//
	clientMock.AssertCalled(t, "Connect")
	clientMock.AssertCalled(t, "Close")
	//
	geoDaoMock.AssertCalled(t, "InsertMany", []model.GeoLocation{{
		IP: "192.168.1.3", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}})

	geoDaoMock.AssertCalled(t, "InsertMany", []model.GeoLocation{{
		IP: "192.168.1.4", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}})

	geoDaoMock.AssertCalled(t, "InsertMany", []model.GeoLocation{{
		IP: "192.168.1.5", CountryCode: "AR", Country: "Iran", City: "Tehran", Latitude: 12, Longitude: 20.236, MysteryValue: "AAAA",
	}})

	geoDaoMock.AssertNumberOfCalls(t, "InsertMany", 3)
}

func TestImporterImpl_Import_FileNotFound(t *testing.T) {
	geoDaoMock := GeoLocationDaoMock{}
	clientMock := ClientWrapperMock{}
	dispatcher := NewDispatcherImpl(&geoDaoMock)
	importer := NewImporterImpl(&clientMock, dispatcher)
	config.C.FileAddress = "not_exist.csv"
	err := importer.Import()
	//goland:noinspection GoNilness
	assert.Equal(t, err.Error(), "open not_exist.csv: no such file or directory")
}
