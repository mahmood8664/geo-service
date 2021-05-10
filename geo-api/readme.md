# Geolocation API
This application provides REST API for geolocation information of
given IP address. Swagger is available at 
http://148.251.100.210:8080/swagger/index.html. 

## Go Version
The minimum version required is 1.15.2

# Building And Deployment
## dependencies
1. **GeoService module:**
You can find this module at [geo-service
](https://github.com/mahmood8664/findhotel-geo/tree/master/geo-service). 
2. **MongoDb** 
## build docker image
to build application by docker execute the following command in the 
[parent](https://github.com/mahmood8664/findhotel-geo) 
directory:
```bash
docker build -t <Image-Name>
```
**Important: In building and deploying process geo-service directory 
must be located beside geo-api directory.** 
## Deployment
If you want to deploy both application and mongodb use below command: 
```bash
docker-compose up --build -d
```
If you want just deploy application, after building and creating docker image you can use 
the following docker command:

```bash
docker run --name container-name --network network-name -p 8080:8080 -d geo-api:latest
```
After deployment successfully, the swagger would be available at 
http://base-url/swagget/index.html

