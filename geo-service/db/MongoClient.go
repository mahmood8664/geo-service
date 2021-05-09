package mongodb

import (
	"context"
	"findhotel.com/geo-service/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	GeoDB         = "geo_db"
	GeoCollection = "geo_locations"
)

var client ClientWrapperImpl

//ClientWrapper is MongoDB Client Wrapper class
type ClientWrapper interface {
	//Connect try to connect to database and will panic if it would not successful
	Connect() error
	//Close try to close database connection and free resources
	Close()
	//Client is original mongo.Client
	Client() *mongo.Client
}

type ClientWrapperImpl struct {
	client     *mongo.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
}

//NewClientWrapperImpl create mongo ClientWrapper
func NewClientWrapperImpl() *ClientWrapperImpl {
	return &client
}

func (r *ClientWrapperImpl) Client() *mongo.Client {
	return r.client
}

func (r *ClientWrapperImpl) Connect() error {
	clientOptions := options.Client().ApplyURI(config.C.MongodbUrl)
	if len(config.C.MongodbUsername) > 0 && len(config.C.MongodbPassword) > 0 {
		clientOptions.SetAuth(options.Credential{
			Username: config.C.MongodbUsername,
			Password: config.C.MongodbPassword,
		})
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		cancelFunc()
		log.Error().Err(err).Msg("")
		return err
	}

	err = client.Ping(context.TODO(), &readpref.ReadPref{}) //try to ping to database and check if it is available
	if err != nil {
		cancelFunc()
		log.Error().Err(err).Msg("")
		return err
	}

	r.client = client
	r.ctx = ctx
	r.cancelFunc = cancelFunc
	return nil
}

func (r *ClientWrapperImpl) Close() {
	r.cancelFunc()
	_ = r.client.Disconnect(r.ctx)
}
