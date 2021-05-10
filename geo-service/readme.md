# GeoLocation Importer and Service Provider Library
This module is developed for two purposes:
1. Provide command line service to import CSV file into database.
2. Provide Model and Service to access geolocation data

## Go Version
The minimum version is 1.15.2

## 1. Use as CSV Importer

### CSV Format
1. Headers
   <br/>
   ip_address,country_code,country,city,latitude,longitude,mystery_value
2. validation
   <br/> All values are **REQUIRED** and invalid records would be **discarded**

* **ip_address**: only IP V4 and V6 format supported
* **country_code**: two letter no matter is capital or not
* **country**: country name
* **city**: city name
* **latitude**: float number between 90 and -90
* **longitude**: float number between 180 and -180
* **mystery_value**: just a string
  <br/>

To use this service for importing data from CSV file into database,
first you need build this application. To build this
service use the following command:

```bash
go build
```

The geo-service binary would be created after building application.  
For more information about importing command use the following command:

```bash
./geo-service import -h
```

The available flags for importer are:

```bash
Flags:
-d, --db-url string     MongoDb URL (default "mongodb://localhost:27017")
-f, --file string       CSV file address
-h, --help              help for import
-p, --password string   MongoDB password
-u, --username string   MongoDB username
```

This is an example of Importing command:

```bash
./geo-service import -f csv-file.csv -d mongodb://localhost:27017 -u user -p 123
```
##Cron Job Script
To run this command as cron job in your server, there is a script file with
the name **script.sh**. It downloads file from a URL and try to import it
into database. First be sure the file has permission to execute. If not just
make it executable by this command:
```bash
chmod +x script.sh
```
Then you can edit crontab by this command:
```bash
crontab -e
```
Then add cron expression of executing the script file to crontab.
For example to execute script at 2 AM every night you can add the following
line to crontab:
```bash
0 2 * * * /path/to/script/script.sh >/dev/null 2>&1
```
## 1. Use as Go library
This module provides an interface ro access geolocation data.
This is an example of using GeoService:
```go
	geoService, err := service.NewGeoService(service.Config{
		MongodbUrl:      config.C.MongoDB.URL,
		MongodbUsername: config.C.MongoDB.Username,
		MongodbPassword: config.C.MongoDB.Password,
	})
	if err != nil {
		fmt.Printf("Cannot connect to service: %s", err.Error())
		return err
	}

	location, err := geoService.GetGeoLocation("192.168.1.1")
```