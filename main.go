package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)
const (
	UserName     string = "root"
	Password     string = "test"
	Addr         string = "db"
	Port         int    = 3306
	Database     string = "test"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)
// Binding from JSON
type Data struct {
	Name string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Customer struct {
	ID string
	Username string
}
func main() {
	

	r := gin.Default()
	r.GET("/ping", index)
	r.GET("/ping/:id", show)
	r.POST("/ping", create)
	r.Run(":80") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong!!!!",
	})
}
func show(c *gin.Context) {
	//組合sql連線字串
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	//連接MySQL
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}
	DB.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	DB.SetMaxOpenConns(MaxOpenConns)
	DB.SetMaxIdleConns(MaxIdleConns)

	id := c.Param("id")
	customer := new(Customer)
	row := DB.QueryRow("select id, username from customer where id=?", id)

	if err := row.Scan(&customer.ID, &customer.Username); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("customer:%+v:", customer)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong!!!!",
		"name": customer.Username,
	})
}


func create(c *gin.Context) {
	var json Data
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//組合sql連線字串
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	//連接MySQL
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}
	DB.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	DB.SetMaxOpenConns(MaxOpenConns)
	DB.SetMaxIdleConns(MaxIdleConns)

	name := json.Name
	password := json.Password
	result, err := DB.Exec("insert INTO customer (username,password) values(?,?)", name, password)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert data id:", lastInsertID)

	c.JSON(http.StatusOK, gin.H{
		"message": "insert success!",
		"id": lastInsertID,
	})
}