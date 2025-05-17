package controllers

import (
	"go-laundry-app/models"
	"go-laundry-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	usecase usecase.TransactionUsecase
	rg *gin.RouterGroup
}

func (t *TransactionController) Route() {
	t.rg.POST("/transactions", t.createTransaction)
	t.rg.GET("/transactions", t.getAllTransaction)
	t.rg.GET("/transactions/:id_bill", t.getTransactionByID)
	t.rg.PUT("/transactions", t.updateTransactions)
	t.rg.DELETE("/transactions/:id_bill", t.deleteTransactionByID)
}

func (t *TransactionController) createTransaction(ctx *gin.Context) {
	var payload models.Transaction

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := t.usecase.CreateNewTransaction(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create Transaction!"})
		return
	}

	completeTransaction, err := t.usecase.GetTransactionByID(transaction.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created transaction"})
		return
	}

	ctx.JSON(http.StatusCreated, completeTransaction)
}

func (t *TransactionController) getAllTransaction(ctx *gin.Context) {
	transactions, err := t.usecase.GetAllTransaction()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Transactions Data!"})
		return
	}

	if len(transactions) > 0 {
		ctx.JSON(http.StatusOK, transactions)
		return
	}
	ctx.JSON(http.StatusOK, transactions)
}

func (t *TransactionController) getTransactionByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id_bill"))

	transactions, err := t.usecase.GetTransactionByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Transaction Not Found"})
		return
	}
	
	ctx.JSON(http.StatusOK, transactions)
}

func (t *TransactionController) updateTransactions(ctx *gin.Context) {
	var payload models.Transaction

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ , err = t.usecase.UpdateTransactionByID(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedTransaction, err := t.usecase.GetTransactionByID(payload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated transaction"})
		return
	}

	ctx.JSON(http.StatusOK, updatedTransaction)
}

func (t *TransactionController) deleteTransactionByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id_bill"))

	err := t.usecase.DeleteTransactionByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func NewTransactionController(usecase usecase.TransactionUsecase, rg *gin.RouterGroup) *TransactionController {
	return &TransactionController{usecase: usecase, rg: rg}
}