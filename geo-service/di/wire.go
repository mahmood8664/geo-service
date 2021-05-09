//+build wireinject

package di

import (
	"findhotel.com/geo-service/csv"
	mongodb "findhotel.com/geo-service/db"
	"github.com/google/wire"
)

func CreateImporter() csv.Importer {
	panic(wire.Build(
		csv.NewImporterImpl,
		CreateDBClient,
		CreateDispatcher,
		wire.Bind(new(csv.Importer), new(csv.ImporterImpl)),
	))
}

func CreateDispatcher() csv.Dispatcher {
	panic(wire.Build(
		csv.NewDispatcherImpl,
		CreateGeoLocationDao,
		wire.Bind(new(csv.Dispatcher), new(*csv.DispatcherImpl)),
	))
}

func CreateGeoLocationDao() mongodb.GeoLocationDao {
	panic(wire.Build(
		mongodb.NewGeoLocationDaoImpl,
		CreateDBClient,
		wire.Bind(new(mongodb.GeoLocationDao), new(mongodb.GeoLocationDaoImpl)),
	))
}

var Set = wire.NewSet(
	//provideMyFooer,
	wire.Bind(new(mongodb.ClientWrapper), new(*mongodb.ClientWrapperImpl)),
)

func CreateDBClient() mongodb.ClientWrapper {
	panic(wire.Build(
		mongodb.NewClientWrapperImpl,
		wire.Bind(new(mongodb.ClientWrapper), new(*mongodb.ClientWrapperImpl)),
	))
}
