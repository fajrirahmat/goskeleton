package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//for real project, don't use somethi like this

//SecretKey just secret key.. better use private key
var SecretKey = "shoulduseprivatekey"

//GetJWTMiddleware return jwt middleware
func GetJWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWT([]byte(SecretKey))
}
