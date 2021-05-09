# Geolocation API
This application provides REST API for geolocation information of
given IP address. Swagger is available at 
http://mamiri.me/goe/swagger/index.html. 

## Go Version
The minimum version is 1.15.2

# Building And Deployment
## dependency
1. GeoService Library <br/>
To deploy this app you need [findhotel-geo-service
](https://github.com/mahmood8664/findhotel-geo-service) 
repository to put at the same directory, so, before building the project 
be sure [findhotel-geo-service
](https://github.com/mahmood8664/findhotel-geo-service) is present.
2. MongoDb 
## build docker image
to build application docker image execute the following command:
```bash
docker build -t <Image-Name>
```
## Deployment
To deploy application in your server run the below command.  
```bash
docker-compose up -d
```
The aforementioned command will deploy both the application and mongodb 
database which is required for this application. 
Application will be started at port 8080 and swagger would be available at 
http://base-url/swagget/index.html

