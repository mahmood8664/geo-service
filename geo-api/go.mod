module findhotel.net/geo-api

go 1.15

replace findhotel.com/geo-service => ../geo-service

require (
	findhotel.com/geo-service v0.0.0-00010101000000-000000000000
	github.com/labstack/echo/v4 v4.2.2
	github.com/mitchellh/mapstructure v1.1.2
	github.com/rs/zerolog v1.21.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
	github.com/ziflex/lecho/v2 v2.3.0
)
