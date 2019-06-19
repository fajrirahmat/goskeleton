package controllers

import (
	"goskeleton/middlewares"
	"goskeleton/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//LoginController struct
type LoginController struct {
	db *gorm.DB
}

func (l *LoginController) Login(ctx echo.Context) error {
	loginRequest := new(LoginRequest)
	if err := ctx.Bind(loginRequest); err != nil {
		return err
	}

	//query user to db
	var user models.User
	l.db.Where("email = ? AND password = ?", loginRequest.Email, loginRequest.Password).Find(&user)
	if user.Email == "" {
		return ctx.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Email and password is invalid",
		})
	}
	//create token
	token := jwt.New(jwt.SigningMethodHS256)

	//set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	t, err := token.SignedString([]byte(middlewares.SecretKey))

	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

//InitializeRoute init route for this Controller
func (l *LoginController) InitializeRoute(e *echo.Echo) {
	e.POST("/login", l.Login)
}

func (l *LoginController) SetDB(db *gorm.DB) {
	l.db = db
}

type (
	//LoginRequest schema for login
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
