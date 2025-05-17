package main

import (
	"fmt"
	"go-laundry-app/config"
	"go-laundry-app/controllers"
	"go-laundry-app/repositories"
	"go-laundry-app/usecase"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	customerUC usecase.CustomerUsecase
	employeeUC usecase.EmployeeUsecase
	productUC usecase.ProductUsecase
	transactionUC usecase.TransactionUsecase
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	controllers.NewCustomerController(s.customerUC, rg).Route()
	controllers.NewEmployeeController(s.employeeUC, rg).Route()
	controllers.NewProductController(s.productUC, rg).Route()
	controllers.NewTransactionController(s.transactionUC, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()

	err := s.engine.Run(s.host)
	if err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		panic("Error Connect to Database")
	}

	customerRepo := repositories.NewCustomerRepository(db)
	customerUseCase := usecase.NewCustomerUseCase(customerRepo)

	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepo)

	productRepo := repositories.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		customerUC: customerUseCase,
		employeeUC: employeeUseCase,
		productUC: productUsecase,
		transactionUC: transactionUsecase,
		engine:     engine,
		host:       host,
	}
}
