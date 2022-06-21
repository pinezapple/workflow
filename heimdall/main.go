// main pkg to start the app
package main

import (
	"github.com/vfluxus/heimdall/boot"
)

// @title 			Heimdall API
// @version 		1.0
// @description 	This is a Heimdall server.
// @termsOfService 	http://swagger.io/terms/
// @contact.name 	API Support
// @contact.url 	http://www.swagger.io/support
// @contact.email 	support@swagger.io
// @license.name 	VinBigData
// @license.url 	http://vinbigdata.org
// @BasePath 		/
// WebServer booting web server by configuration
func main() {
	boot.Boot()
}
