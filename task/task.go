package task

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	os.Setenv("DATABASE_URL", "postgres://izwjgrts:wHgZyqD67n2-Bf9Y3sCgD0JKK_rKhGXM@lallah.db.elephantsql.com:5432/izwjgrts")
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
}

func Create_table() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT,
	status TEXT
	);
	`
	_, err := db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	fmt.Println("create table success")
}

func CreateCustomersHandler(c *gin.Context) {
	cus := Customer{}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	row := db.QueryRow("INSERT INTO customers (name,email,status) values ($1, $2, $3) RETURNING id", cus.Name, cus.Email, cus.Status)
	err := row.Scan(&cus.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, cus)
}

func GetCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cus := &Customer{}

	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, cus)
}

func GetCustomersHandler(c *gin.Context) {
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	customers := []Customer{}
	for rows.Next() {
		cus := Customer{}
		err := rows.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		customers = append(customers, cus)
	}
	c.JSON(http.StatusOK, customers)
}

func UpdateCustomerHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	row := stmt.QueryRow(id)

	cus := &Customer{}
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.ShouldBindJSON(cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if _, err := stmt.Exec(id, cus.Name, cus.Email, cus.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, cus)
}

func DeleteCustomersHandler(c *gin.Context) {
	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
