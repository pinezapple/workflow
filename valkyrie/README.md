# valkyrie
Vfluxus file manager service.


## Development
Remember to write OpenAPI document.  
User tool swag to write document foreach endpoint

### Tutorial 
Github url: https://github.com/swaggo/swag

For echo framework install swag
```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/echo-swagger
```
(Remember to add GOPATH/bin to your PATH environment)  


Follow instruction in echo swagger middleware to add comment for each endpoint: https://github.com/swaggo/echo-swagger
After comment on all endpoint, please run the command
```
swag init 
```
to generate the OpenAPI docs and a Go module.   
Import the generated Go module to source code
```Go
import _ "github.com/vfluxus/valkyrie/docs
```
and add the to the Router 
```
  import "github.com/swaggo/echo-swagger"
  ...

  e := echo.New()
	e.GET("/docs/*", echoSwagger.WrapHandler)

```
Run the command 
```
go run . 
```
and browse to url http://localhost:10006/docs/index.html to see the document.

