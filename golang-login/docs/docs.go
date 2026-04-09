package docs

import "github.com/swaggo/swag"

var SwaggerInfo = &swag.Spec{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/",
	Title:       "Golang Login API",
	Description: "API de autenticacao",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
