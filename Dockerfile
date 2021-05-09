FROM golang:1.15.2-alpine3.12 AS build
RUN mkdir app
COPY . /app
COPY ../findhotel-geo-service /geo-service
WORKDIR /app
RUN go build && \
    mkdir /out && \
    cp ./geo-api /out/ && \
    cp -r /app/resources/ /out/
FROM alpine:3.12
RUN mkdir /app
WORKDIR /app
COPY --from=build /out /app
EXPOSE 8080