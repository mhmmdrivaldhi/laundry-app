package controllers

import (
	"go-laundry-app/models"
	"go-laundry-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	useCase usecase.CustomerUsecase
	rg      *gin.RouterGroup
}

func (c *CustomerController) Route() {
	c.rg.POST("/customers", c.createNewCustomer)
	c.rg.GET("/customers", c.getAllCustomer)
	c.rg.GET("/customers/:id", c.getCustomerById)
	c.rg.PUT("/customers", c.updateCustomerById)
	c.rg.DELETE("/customers/:id", c.deleteCustomerById)
}

func (c *CustomerController) createNewCustomer(ctx *gin.Context) {
	var payload models.Customer

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := c.useCase.CreateNewCustomer(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create Data Customer!"})
		return
	}
	ctx.JSON(http.StatusCreated, customer)
}

func (c *CustomerController) getAllCustomer(ctx *gin.Context) {
	customers, err := c.useCase.GetAllCustomer()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Customers Data!"})
		return
	}
	if len(customers) > 0 {
		ctx.JSON(http.StatusOK, customers)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List customers is empty"})
}

func (c *CustomerController) getCustomerById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	customer, err := c.useCase.GetCustomerById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer by ID"})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) updateCustomerById(ctx *gin.Context) {
	var payload models.Customer

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := c.useCase.UpdateCustomerById(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) deleteCustomerById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.useCase.DeleteCustomerById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func NewCustomerController(useCase usecase.CustomerUsecase, rg *gin.RouterGroup) *CustomerController {
	return &CustomerController{useCase: useCase, rg: rg}
}
