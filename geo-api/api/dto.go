package api

import "findhotel.com/geo-service/model"

type GeoLocationResponse struct {
	IP           string  `json:"ip"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MysteryValue string  `json:"mystery_value"`
}

//FromGeoLocationModel convert model.GeoLocation to GeoLocationResponse
func FromGeoLocationModel(geoLocation model.GeoLocation) GeoLocationResponse {
	return GeoLocationResponse{
		IP:           geoLocation.IP,
		City:         geoLocation.City,
		MysteryValue: geoLocation.City,
		Latitude:     geoLocation.Latitude,
		Longitude:    geoLocation.Longitude,
		Country:      geoLocation.Country,
		CountryCode:  geoLocation.CountryCode,
	}
}
