package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Controller interface {
	InitializeRoute(e *echo.Echo)
	SetDB(db *gorm.DB)
}
