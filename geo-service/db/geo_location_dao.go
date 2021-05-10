package mongodb

import (
	"context"
	"findhotel.com/geo-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

//GeoLocationDao is a interface to perform database operation on
type GeoLocationDao interface {
	//Insert a model.GeoLocation and return inserter object id (hex string) or error
	Insert(geoLocation model.GeoLocation) (id string, err error)
	//GetOne try to find model.GeoLocation by ip address or return error
	GetOne(ip string) (user model.GeoLocation, err error)
	//InsertMany try to insert slice of model.GeoLocation(s), It is not act transactional, so both ids and err can have
	//valid value, ids represents successful inserted Ids and err (if it is kind of mongo.BulkWriteException) contains
	//indexes which failed to insert into database
	InsertMany(geoLocation []model.GeoLocation) (ids []interface{}, err error)
	//CreateIndex create necessary indexes of database collections
	CreateIndex() (err error)
}

type GeoLocationDaoImpl struct {
	clientWrapper ClientWrapper
	geoCollection *mongo.Collection
}

func NewGeoLocationDaoImpl(wrapper ClientWrapper) GeoLocationDaoImpl {
	return GeoLocationDaoImpl{
		clientWrapper: wrapper,
	}
}

func (r GeoLocationDaoImpl) GetOne(ip string) (geoLocation model.GeoLocation, err error) {
	result := r.GeoLocationConnection().FindOne(context.TODO(), bson.D{{"ip", ip}})
	err = result.Decode(&geoLocation)
	return geoLocation, err
}

func (r GeoLocationDaoImpl) Insert(geoLocation model.GeoLocation) (id string, err error) {
	one, err := r.GeoLocationConnection().InsertOne(context.TODO(), geoLocation)
	if err != nil {
		return "", err
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r GeoLocationDaoImpl) InsertMany(geoLocation []model.GeoLocation) (ids []interface{}, err error) {
	var geos = make([]interface{}, len(geoLocation))
	for i, g := range geoLocation {
		geos[i] = g
	}
	ordered := false
	many, err := r.GeoLocationConnection().InsertMany(context.TODO(), geos, &options.InsertManyOptions{Ordered: &ordered})
	if err != nil {
		if many != nil {
			return many.InsertedIDs, err
		} else {
			return nil, err
		}
	}
	return many.InsertedIDs, nil
}

func (r GeoLocationDaoImpl) CreateIndex() (err error) {
	_, err = r.GeoLocationConnection().Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"ip", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		})
	return err
}

func (r GeoLocationDaoImpl) GeoLocationConnection() *mongo.Collection {
	if r.geoCollection == nil {
		r.geoCollection = r.clientWrapper.Client().Database(GeoDB).Collection(GeoCollection)
	}
	return r.geoCollection
}
