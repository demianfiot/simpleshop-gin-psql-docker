package handler

import (
	"net/http"
	"prac/todo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// DTO  zminyty
	var input struct {
		Items []struct {
			ProductID uint    `json:"product_id" binding:"required"`
			Quantity  int     `json:"quantity" binding:"required,min=1"`
			Price     float64 `json:"price" binding:"required,min=0"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// DTO  zminyty
	var orderItems []todo.OrderItem
	for _, item := range input.Items {
		orderItems = append(orderItems, todo.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}
	ctx := c.Request.Context()
	orderID, err := h.services.Order.CreateOrder(ctx, userID, orderItems)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{ // 201 Created
		"order_id": orderID,
		"message":  "Order created successfully",
	})
}

func (h *Handler) GetUserOrders(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	ctx := c.Request.Context()
	orders, err := h.services.Order.GetUserOrders(ctx, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (h *Handler) GetAllOrders(c *gin.Context) {
	// tilky admin - middleware
	ctx := c.Request.Context()
	orders, err := h.services.Order.GetAllOrders(ctx)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (h *Handler) GetOrderByID(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid order id param")
		return
	}
	ctx := c.Request.Context()
	order, err := h.services.Order.GetOrderByID(ctx, uint(orderID))
	if err != nil {
		if err.Error() == "order not found" {
			NewErrorResponse(c, http.StatusNotFound, "order not found")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// prava dostupu
	currentUserID, _ := h.getUserID(c)
	currentUserRole, _ := c.Get("userRole")

	if currentUserRole != "admin" && order.UserID != currentUserID {
		NewErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid order id param")
		return
	}

	// DTO тstatusu (sminyty)
	var input struct {
		Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
	}

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	err = h.services.Order.UpdateOrderStatus(ctx, uint(orderID), input.Status)
	if err != nil {
		if err.Error() == "order not found" {
			NewErrorResponse(c, http.StatusNotFound, "order not found")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}
