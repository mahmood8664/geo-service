package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"regexp"
	"strconv"
	"strings"
)

//GeoLocation is a model to keep geolocation data
type GeoLocation struct {
	Id           *primitive.ObjectID `bson:"_id,omitempty" json:""`
	IP           string              `bson:"ip" json:"ip"`
	CountryCode  string              `bson:"country_code" json:"country_code"`
	Country      string              `bson:"country" json:"country"`
	City         string              `bson:"city" json:"city"`
	Latitude     float64             `bson:"latitude" json:"latitude"`
	Longitude    float64             `bson:"longitude" json:"longitude"`
	MysteryValue string              `bson:"mystery_value" json:"mystery_value"`
}

var (
	parseError          = errors.New("error during parsing csv record")
	countryCodeRegex, _ = regexp.Compile("[a-zA-Z]{2}")
)

//ParseGeoLocation parse a slice of string and tries to convert it into GeoLocation, if error would not nil it means
// conversion was not successful and GeoLocation is not valid
func ParseGeoLocation(record []string) (GeoLocation, error) {
	var geoLocation = GeoLocation{}

	if len(record) != 7 {
		return geoLocation, parseError
	}

	for _, s := range record {
		if len(strings.Trim(s, " \"'")) == 0 {
			return geoLocation, parseError
		}
	}

	ip := net.ParseIP(strings.Trim(record[0], " '\""))
	if ip == nil {
		return geoLocation, parseError
	}

	geoLocation.IP = ip.String()

	if !countryCodeRegex.Match([]byte(record[1])) || len(strings.Trim(record[1], " '\"")) != 2 {
		return geoLocation, parseError
	}
	geoLocation.CountryCode = record[1]

	geoLocation.Country = strings.Trim(record[2], " '\"")
	geoLocation.City = strings.Trim(record[3], " '\"")

	lat, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return geoLocation, parseError
	}
	if lat >= -90 && lat <= 90 {
		geoLocation.Latitude = lat
	} else {
		return geoLocation, parseError
	}

	long, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return geoLocation, parseError
	}
	if long >= -180 && long <= 180 {
		geoLocation.Longitude = long
	} else {
		return geoLocation, parseError
	}

	geoLocation.MysteryValue = record[6]

	return geoLocation, nil
}
