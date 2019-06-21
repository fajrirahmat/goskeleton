package controllers

import (
	"goskeleton/middlewares"
	"goskeleton/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type ProductController struct {
	db *gorm.DB
}

func (p *ProductController) List(ctx echo.Context) error {
	var products []models.Product
	p.db.Find(&products)
	return ctx.JSON(http.StatusOK, products)
}

func (p *ProductController) Add(ctx echo.Context) error {
	addPrdRequest := new(AddProductRequest)
	if err := ctx.Bind(addPrdRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to parse request",
		})
	}

	p.db.Create(&models.Product{
		Name:   addPrdRequest.Name,
		Price:  addPrdRequest.Price,
		Stock:  addPrdRequest.Stock,
		UserID: p.getUserIDfromToken(ctx),
	})

	return ctx.JSON(http.StatusCreated, map[string]string{
		"message": "Product created",
	})
}

func (p *ProductController) getUserIDfromToken(ctx echo.Context) uint {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	emailFromToken := claims["email"].(string)

	var userDB models.User
	p.db.Find(&userDB).Where("email = ?", emailFromToken)
	return userDB.ID
}

//InitializeRoute init route for this Controller
func (p *ProductController) InitializeRoute(e *echo.Echo) {
	r := e.Group("/products", middlewares.GetJWTMiddleware())
	r.GET("", p.List)
	r.POST("/add", p.Add)
}

//SetDB inject DB
func (p *ProductController) SetDB(db *gorm.DB) {
	p.db = db
}

type (
	//AddProductRequest request to add product
	AddProductRequest struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Stock int     `json:"stock"`
	}
)
