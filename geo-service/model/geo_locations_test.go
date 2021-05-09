package model

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseGeoLocation(t *testing.T) {

	records := MockData()

	geoLocations := make([]GeoLocation, 0)
	for _, record := range records {
		location, err := ParseGeoLocation(record)
		if err != nil {
			log.Info().Err(err)
		} else {
			geoLocations = append(geoLocations, location)
		}
	}

	assert.Equal(t, len(geoLocations), 24)
	assert.Equal(t, geoLocations[0], GeoLocation{"", "255.255.255.255", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[1], GeoLocation{"", "1.1.1.1", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[2], GeoLocation{"", "2.2.2.2", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[3], GeoLocation{"", "0.0.0.0", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[4], GeoLocation{"", "192.168.1.3", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[5], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[6], GeoLocation{"", "192.168.1.1", "AR", "I", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[7], GeoLocation{"", "192.168.1.1", "AR", "A B", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[8], GeoLocation{"", "192.168.1.1", "AR", "AA BB", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[9], GeoLocation{"", "192.168.1.1", "AR", "AAA BBB", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[10], GeoLocation{"", "192.168.1.1", "AR", "AAAA BBBB", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[11], GeoLocation{"", "192.168.1.1", "AR", "AAAAA BBBBB", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[12], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[13], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Paris", 12, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[14], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 0, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[15], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 90, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[16], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", -90, 20.236, "AAAA"})
	assert.Equal(t, geoLocations[17], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, 0, "AAAA"})
	assert.Equal(t, geoLocations[18], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, 180, "AAAA"})
	assert.Equal(t, geoLocations[19], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, -180, "AAAA"})
	assert.Equal(t, geoLocations[20], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, 25, "a"})
	assert.Equal(t, geoLocations[21], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, 25, "aa"})
	assert.Equal(t, geoLocations[22], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, 25, "'aaa'"})
	assert.Equal(t, geoLocations[23], GeoLocation{"", "192.168.1.1", "AR", "Iran", "Tehran", 10, 25, "'aaaa"})
}

func BenchmarkParseGeoLocation(b *testing.B) {
	records := MockData()
	for i := 0; i < b.N; i++ {
		geoLocations := make([]GeoLocation, 0)
		for _, record := range records {
			location, err := ParseGeoLocation(record)
			if err != nil {
				log.Info().Err(err)
			} else {
				geoLocations = append(geoLocations, location)
			}
		}
	}
}

func MockData() [][]string {
	return [][]string{
		{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"},
		{"", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{" ", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"  ", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"235", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"235.325.36.365", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"255.32.63", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"255.255.255.255", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"}, //OK
		{"1.1.1.1", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},         //OK
		{" 2.2.2.2 ", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},       //OK
		{"0.0.0.0", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},         //OK
		{"192.168.1.003", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"},   //OK
		{"192.168.1.1", "", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", " ", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "  ", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "A", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.1", "ARR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "ARRR", "Iran", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", " ", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "  ", "Tehran", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "I", "Tehran", "12", "20.236", "AAAA"},               //OK
		{"192.168.1.1", "AR", "A B", "Tehran", "12", "20.236", "AAAA"},             //OK
		{"192.168.1.1", "AR", "'AA BB'", "Tehran", "12", "20.236", "AAAA"},         //OK
		{"192.168.1.1", "AR", "'AAA BBB", "Tehran", "12", "20.236", "AAAA"},        //OK
		{"192.168.1.1", "AR", "\"AAAA BBBB\"", "Tehran", "12", "20.236", "AAAA"},   //OK
		{"192.168.1.1", "AR", "  AAAAA BBBBB  ", "Tehran", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.1", "AR", "Iran", "", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", " ", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "  ", "12", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "  Tehran  ", "12", "20.236", "AAAA"}, //OK
		{"192.168.1.1", "AR", "Iran", "'Paris'", "12", "20.236", "AAAA"},    //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "0", "20.236", "AAAA"},      //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "91", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "-91", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "90", "20.236", "AAAA"},  //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "-90", "20.236", "AAAA"}, //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "-90.000001", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "90.000001", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", " ", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "  ", "20.236", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", " ", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "  ", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "0", "AAAA"},    //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "180", "AAAA"},  //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "-180", "AAAA"}, //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "180.0001", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "-180.0001", "AAAA"},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", ""},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", " "},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", "  "},
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", "a"},     //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", "aa"},    //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", "'aaa'"}, //OK
		{"192.168.1.1", "AR", "Iran", "Tehran", "10", "25", "'aaaa"}, //OK
		{"", "", "", "", "", "", ""},
	}
}
