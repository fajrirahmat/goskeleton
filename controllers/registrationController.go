package controllers

import (
	"errors"
	"goskeleton/models"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//RegistrationController controller struct for registration process
type RegistrationController struct {
	db *gorm.DB
}

//Register handler for registration process
func (r *RegistrationController) Register(ctx echo.Context) error {
	registrationRequest := new(RegistrationRequest)
	if err := ctx.Bind(registrationRequest); err != nil {
		return err
	}

	var existingUser models.User
	r.db.First(&existingUser, "email = ? ", registrationRequest.Email)
	if existingUser.Email != "" {
		return errors.New("User with the email has been registered")
	}

	r.db.Create(&models.User{
		Email:    registrationRequest.Email,
		Name:     registrationRequest.Name,
		Password: registrationRequest.Password,
	})

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User successfully created",
	})

}

//SetDB inject DB
func (r *RegistrationController) SetDB(db *gorm.DB) {
	r.db = db
}

//InitializeRoute RegistrationController route
func (r *RegistrationController) InitializeRoute(e *echo.Echo) {
	e.POST("/registration", r.Register)
}

type (
	//RegistrationRequest for request body
	RegistrationRequest struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
)
