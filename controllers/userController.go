package controllers

import (
	"goskeleton/middlewares"
	"goskeleton/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserController struct {
	db *gorm.DB
}

func (u *UserController) GetProfile(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	emailFromToken := claims["email"].(string)

	email := ctx.Param("email")
	if emailFromToken != email {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "This profile is not belong to you",
		})
	}

	var profile models.User
	u.db.Where("email = ?", email).Find(&profile)
	profile.Password = ""
	return ctx.JSON(http.StatusOK, profile)
}

func (u *UserController) GetUserOrderHistory(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	emailFromToken := claims["email"].(string)

	email := ctx.Param("email")
	if emailFromToken != email {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "This profile is not belong to you",
		})
	}

	var profile models.User
	u.db.Where("email = ?", email).Find(&profile)
	var orders []models.Order
	u.db.Where("user_id = ?", profile.ID).Find(&orders)
	for i := range orders {
		var p models.Product
		o := orders[i]
		u.db.Model(&o).Related(&p)
		p.Stock = 0
		orders[i].Product = p
	}
	return ctx.JSON(http.StatusOK, orders)
}

//InitializeRoute init route for this Controller
func (u *UserController) InitializeRoute(e *echo.Echo) {
	e.GET("/users/:email", u.GetProfile, middlewares.GetJWTMiddleware())
	e.GET("/users/:email/history", u.GetUserOrderHistory, middlewares.GetJWTMiddleware())
}

//SetDB inject DB
func (u *UserController) SetDB(db *gorm.DB) {
	u.db = db
}
