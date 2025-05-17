package controllers

import (
	"go-laundry-app/models"
	"go-laundry-app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	usecase usecase.EmployeeUsecase
	rg *gin.RouterGroup
}

func (e *EmployeeController) Route() {
	e.rg.POST("/employees", e.createNewEmployee)
	e.rg.GET("/employees", e.getAllEmployee)
	e.rg.GET("/employees/:id", e.getEmployeeByID)
	e.rg.PUT("/employees", e.updateEmployeeByID)
	e.rg.DELETE("/employees/:id", e.deleteEmployeeByID)
}

func (e *EmployeeController) createNewEmployee(ctx *gin.Context) {
	var payload models.Employee

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := e.usecase.CreateNewEmployee(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create Data Employee!"})
		return
	}
	ctx.JSON(http.StatusCreated, employee)
}

func (e *EmployeeController) getAllEmployee(ctx *gin.Context) {
	employees, err := e.usecase.GetAllEmployee()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Employees Data!"})
		return
	}
	if len(employees) > 0 {
		ctx.JSON(http.StatusOK, employees)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List employees is empty"})
}

func (e *EmployeeController) getEmployeeByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	employee, err := e.usecase.GetEmployeeById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get employee by ID"})
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (e *EmployeeController) updateEmployeeByID(ctx *gin.Context) {
	var payload models.Employee

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := e.usecase.UpdateEmployeeById(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

func (e *EmployeeController) deleteEmployeeByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := e.usecase.DeleteEmployeeById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func NewEmployeeController(usecase usecase.EmployeeUsecase, rg *gin.RouterGroup) *EmployeeController {
	return &EmployeeController{usecase: usecase, rg: rg}
}