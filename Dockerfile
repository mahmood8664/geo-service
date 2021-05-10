#this is Dockerfile for building geo-api
#Have to be executed at parent directory and both geo-api and geo-service module must be present
FROM golang:1.15.2-alpine3.12 AS build
RUN mkdir app
COPY . /app
WORKDIR /app/geo-api
RUN go build && \
    mkdir /out && \
    cp ./geo-api /out/ && \
    cp -r ./resources/ /out/
FROM alpine:3.12
RUN mkdir /app
WORKDIR /app
COPY --from=build /out /app
EXPOSE 8080
CMD ["./geo-api", "start"]