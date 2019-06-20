package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

//HelloController hello controller struct
type HelloController struct{}

//HelloMessage helo message struct
type HelloMessage struct {
	Message string `json:"message"`
}

//HelloWorld function to handle request to return Hello world
func (h *HelloController) HelloWorld(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, &HelloMessage{
		Message: "Hello World",
	})
}

//InitializeRoute implementation from InitializeRoute of BaseController
func (h *HelloController) InitializeRoute(e *echo.Echo) {
	var SecretKey = "shoulduseprivatekey"
	e.GET("/", h.HelloWorld, middleware.JWT([]byte(SecretKey)))
}

//SetDB no db to set, but because register
func (h *HelloController) SetDB(db *gorm.DB) {}
