package handler

import (
	"net/http"
	"prac/todo"
	"strconv"

	"prac/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	var input todo.CreateProductInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// check ID potochngo korustuvacha
	sellerID, err := h.getUserID(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := h.services.Product.CreateProduct(input, sellerID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := h.services.Product.GetAllProducts()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func (h *Handler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid product id param")
		return
	}

	product, err := h.services.Product.GetProductByID(uint(id))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid product id param")
		return
	}

	var input todo.UpdateProductInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// check ID dl9 prav
	currentUserID, err := h.getUserID(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	updatedProduct, err := h.services.Product.UpdateProduct(uint(productID), input, currentUserID)
	if err != nil {
		if err == repository.ErrProductNotFound {
			NewErrorResponse(c, http.StatusNotFound, "product not found")
			return
		}
		if err == repository.ErrAccessDenied {
			NewErrorResponse(c, http.StatusForbidden, "access denied")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": updatedProduct,
		"message": "Product updated successfully",
	})
}
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}
	err = h.services.Product.DeleteProduct(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
