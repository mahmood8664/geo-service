# Geolocation API
This application provides REST API for geolocation information of
given IP address. Swagger is available at 
http://148.251.100.210:8080/swagger/index.html. 

## Go Version
The minimum version is 1.15.2

# Building And Deployment
## dependencies
1. **GeoService module:**
You can find this module at [geo-service
](https://github.com/mahmood8664/findhotel-geo/tree/master/geo-service). 
2. **MongoDb** 
## build docker image
to build application by docker execute the following command:
```bash
docker build -t <Image-Name>
```
**Important: In building process geo-service directory must be located 
beside geo-api directory.** 
## Deployment
To deploy application in your server run the below command.   
```bash
docker-compose up --build -d
```
The aforementioned command will deploy both the application and mongodb 
database which is required for this application. 
Application exposes port 8080 and swagger would be available at 
http://base-url/swagget/index.html

