package controllers

import (
	"goskeleton/middlewares"
	"goskeleton/models"
	"goskeleton/status"
	"net/http"
	"strconv"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type OrderController struct {
	db *gorm.DB
}

var mutex = &sync.Mutex{}

func (o *OrderController) Create(ctx echo.Context) error {
	orderRequest := new(CreateOrderRequest)
	if err := ctx.Bind(orderRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if orderRequest.Qty == 0 {
		orderRequest.Qty = 1
	}
	tx := o.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var product models.Product
	tx.Where("id = ?", orderRequest.ProductID).Find(&product)
	if (product.Stock - orderRequest.Qty) < 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Out of stock",
		})
	}
	product.Stock = product.Stock - orderRequest.Qty
	order := &models.Order{
		Product:   product,
		ProductID: product.ID,
		Qty:       orderRequest.Qty,
		Status:    status.New,
		UserID:    o.getUserIDfromToken(ctx),
	}
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"orderId": order.ID,
	})
}

func (o *OrderController) getUserIDfromToken(ctx echo.Context) uint {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	emailFromToken := claims["email"].(string)

	var userDB models.User
	o.db.Find(&userDB).Where("email = ?", emailFromToken)
	return userDB.ID
}

func (o *OrderController) Cancel(ctx echo.Context) error {
	strOrderID := ctx.Param("orderid")
	if strOrderID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Please provided order id information",
		})
	}
	orderID, err := strconv.ParseUint(strOrderID, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	tx := o.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var order models.Order
	var product models.Product
	tx.Where("id = ?", orderID).Find(&order)
	tx.Model(&order).Related(&product)
	//delete only change deleted time
	order.Status = status.Canceled
	if err = tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	product.Stock = product.Stock + order.Qty
	if err = tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Where("id = ?", orderID).Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return ctx.JSON(http.StatusAccepted, map[string]string{
		"message": "Order has been canceled",
	})
}

func (o *OrderController) Update(ctx echo.Context) error {
	updateOrderRequest := new(UpdateOrderRequest)
	if err := ctx.Bind(updateOrderRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	tx := o.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var order models.Order
	tx.Where("id = ?", updateOrderRequest.OrderID).Find(&order)
	var product models.Product
	tx.Model(&order).Related(&product)
	if (product.Stock + order.Qty - updateOrderRequest.Qty) < 0 {
		tx.Rollback()
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Out of stock",
		})
	}
	product.Stock = product.Stock + order.Qty - updateOrderRequest.Qty
	order.Qty = updateOrderRequest.Qty

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return ctx.JSON(http.StatusAccepted, map[string]string{
		"message": "Order is updated",
	})
}

func (o *OrderController) Pay(ctx echo.Context) error {
	payOrderRequest := new(PayOrderRequest)
	if err := ctx.Bind(payOrderRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	tx := o.db.Begin()
	var order models.Order
	if err := tx.Where("id = ?", payOrderRequest.OrderID).Find(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if order.Status != status.New {
		return ctx.JSON(http.StatusExpectationFailed, map[string]string{
			"message": "Order can't be paid",
		})
	}

	//TODO call payment integration here
	order.Status = status.Paid

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return ctx.JSON(http.StatusAccepted, map[string]string{
		"message": "Order has been paid",
	})
}

func (o *OrderController) InitializeRoute(e *echo.Echo) {
	r := e.Group("/purchase", middlewares.GetJWTMiddleware())
	r.POST("/create", o.Create)
	r.DELETE("/cancel/:orderid", o.Cancel)
	r.PUT("/update", o.Update)
	r.POST("/pay", o.Pay)
}
func (o *OrderController) SetDB(db *gorm.DB) {
	o.db = db
}

type (
	CreateOrderRequest struct {
		ProductID uint `json:"productId"`
		Qty       int  `json:"qty"`
	}
	UpdateOrderRequest struct {
		OrderID uint `json:"orderId"`
		Qty     int  `json:"qty"`
	}

	PayOrderRequest struct {
		//assume we have wallet id and our store only can be paid using the wallet... :P
		WalletID uint `json:"walletId"`
		OrderID  uint `json:"orderId"`
	}
)
