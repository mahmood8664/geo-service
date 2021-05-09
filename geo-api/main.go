package main

import (
	"findhotel.net/geo-api/cmd"
	_ "findhotel.net/geo-api/docs"
)

// @title GeoLocation API
// @version 1.0
// @description GeoLocation API
// @termsOfService http://swagger.io/terms/
// @contact.name Mahmoud AllamehAmiri
// @contact.email m.allamehamiri@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	cmd.Execute()
}
