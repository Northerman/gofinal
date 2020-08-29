package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/northerman/gofinal/middleware"
	"github.com/northerman/gofinal/task"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	task.Create_table()
	customers := r.Group("/customers")
	customers.Use(middleware.Auth)
	customers.POST("", task.CreateCustomersHandler)
	customers.GET("/:id", task.GetCustomerByIdHandler)
	customers.GET("", task.GetCustomersHandler)
	customers.PUT("/:id", task.UpdateCustomerHandler)
	customers.DELETE("/:id", task.DeleteCustomersHandler)
	return r
}

func main() {
	fmt.Println("customer service")
	r := setupRouter()
	r.Run(":2009")
	//run port ":2009"
}
