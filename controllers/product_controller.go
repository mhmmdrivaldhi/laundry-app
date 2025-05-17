package controllers

import (
	"go-laundry-app/models"
	"go-laundry-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	usecase usecase.ProductUsecase
	rg *gin.RouterGroup
}

func ( p *ProductController) Route() {
	p.rg.POST("/products", p.createNewProduct)
	p.rg.GET("/products", p.getAllProduct)
	p.rg.GET("/products/:id", p.getProductByID)
	p.rg.PUT("/products", p.updateProductByID)
	p.rg.DELETE("/products/:id", p.deleteProductByID)
}

func (p *ProductController) createNewProduct(ctx *gin.Context) {
	var payload models.Product

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := p.usecase.CreateNewProduct(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Data Product"})
		return
	}
	ctx.JSON(http.StatusCreated, product)
}

func (p *ProductController) getAllProduct(ctx *gin.Context) {
	products, err := p.usecase.GetAllProduct()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Product Daata!"})
		return
	}

	if len(products) > 0 {
		ctx.JSON(http.StatusOK, products)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List Product is empty"})
}

func (p *ProductController) getProductByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	product, err := p.usecase.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product by ID"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (p *ProductController) updateProductByID(ctx *gin.Context) {
	var payload models.Product

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := p.usecase.UpdateProductByID(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (p *ProductController) deleteProductByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := p.usecase.DeleteProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func NewProductController(usecase usecase.ProductUsecase, rg *gin.RouterGroup) *ProductController {
	return &ProductController{usecase: usecase, rg: rg}
}