package main

import (
	"github.com/joho/godotenv"
	// "golang.org/x/net/html/atom"

	//"log"
	"net/http"
	"os"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init(){
	err := godotenv.Load()
	if err!=nil {
		fmt.Println(err)
	}

	dbName := os.Getenv("db_name")
	dbUser := os.Getenv("db_user")
	dbPass := os.Getenv("db_pass")
	dbhost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%v user=%v dbname=%v sslmode=disable password=%v", dbhost, dbUser, dbName, dbPass)
	fmt.Println(dbUri)

	db, err = gorm.Open("postgres", dbUri)
	if err!=nil {
		fmt.Println(err)
	}
	// var td = todoModel{Title: "asdfasdf", Completed: 1}
	// db.Create(&td)
	db.AutoMigrate()
}

type(
	todoModel struct{
		gorm.Model
		Title string `json:"title"`
		Completed int `json:"completed"`
	}

	transformtodoModel struct{
		ID uint
		Title string `json:"title"`
		Completed bool `json:"completed"`
	}
)

func createTodo(c *gin.Context){
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	title :=c.PostForm("title")
	todo := todoModel{Title: title, Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "create todo successfully", "resourceID": todo.ID})
}

func updateTodo(c *gin.Context){
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)

	if todo.ID == 0{
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "not found todo"})
		return
	}

	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	title := c.PostForm("title")
	db.Model(&todo).Update("title", title)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "update todo successfully"})
}

func deleteTodo(c *gin.Context){
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0{
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "not found todo"})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "delete to successfully"})
}

func fetchSingleTodo(c *gin.Context){
	var todo todoModel
	var _todo []transformtodoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)
	fmt.Println("todoID", todo)

	if todo.ID== 0{
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "not found todo"})
		return
	}
	completed := false
	if todo.Completed == 1{
		completed = true
	}else{
		completed = false
	}

	_todo = append(_todo, transformtodoModel{ID: todo.ID, Title: todo.Title, Completed: completed})
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "fetch todo successfully !!!", "data": _todo})
}

func fetchAllTodos(c *gin.Context){
	var todos []todoModel
	var _todos []transformtodoModel

	db.Find(&todos)

	if len(todos) <=0{
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusOK, "message": "failed fetch all todo"})
		return
	}

	for _, items:= range todos {
		completed := false
		if items.Completed == 1{
			completed = true
		}else{
			completed = false
		}
		_todos = append(_todos, transformtodoModel{ID: items.ID, Completed: completed, Title: items. Title})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "fetch all todo successfully", "data": _todos})
}

func main(){
	router := gin.Default()
	v1 := router.Group("api/v1/todos")
	{
		v1.GET("/", fetchAllTodos)
		v1.POST("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
		v1.POST("/", createTodo)
		v1.GET("/:id", fetchSingleTodo)
	}
	port := os.Getenv("port")
	if port == ""{
		port = ":8080"
	}
	router.Run(port)
}