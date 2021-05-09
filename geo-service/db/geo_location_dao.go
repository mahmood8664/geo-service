package mongodb

import (
	"context"
	"findhotel.com/geo-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
}

type GeoLocationDaoImpl struct {
	clientWrapper ClientWrapper
}

func NewGeoLocationDaoImpl(wrapper ClientWrapper) GeoLocationDaoImpl {
	return GeoLocationDaoImpl{
		clientWrapper: wrapper,
	}
}

func (r GeoLocationDaoImpl) GetOne(ip string) (geoLocation model.GeoLocation, err error) {
	result := r.clientWrapper.Client().Database(GeoDB).Collection(GeoCollection).FindOne(context.TODO(), bson.D{{"ip", ip}})
	err = result.Decode(&geoLocation)
	return geoLocation, err
}

func (r GeoLocationDaoImpl) Insert(geoLocation model.GeoLocation) (id string, err error) {
	one, err := r.clientWrapper.Client().Database(GeoDB).Collection(GeoCollection).InsertOne(context.TODO(), geoLocation)
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
	many, err := r.clientWrapper.Client().Database(GeoDB).Collection(GeoCollection).InsertMany(context.TODO(), geos, &options.InsertManyOptions{Ordered: &ordered})
	if err != nil {
		if many != nil {
			return many.InsertedIDs, err
		} else {
			return nil, err
		}
	}
	return many.InsertedIDs, nil
}
