package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_crud")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLFiles("index.html", "templates/create.html", "templates/edit.html", "templates/view.html")

	router.Static("/static", "./assets")

	router.GET("/", index)
	router.GET("/create", create)
	router.POST("/store", store)
	router.GET("/edit", edit)
	router.POST("/update", update)
	router.GET("/view", view)
	router.GET("/delete", delete)

	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}

func index(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}
		users = append(users, user)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Users": users,
	})
}

func create(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func store(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func edit(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID tidak ditemukan",
		})
		return
	}

	var user User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"error": "Pengguna tidak ditemukan",
			})
		} else {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.HTML(http.StatusOK, "edit.html", user)
}

func update(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")

	_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", name, email, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func view(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID tidak ditemukan",
		})
		return
	}

	var user User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"error": "Pengguna tidak ditemukan",
			})
		} else {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.HTML(http.StatusOK, "view.html", user)
}

func delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID tidak ditemukan",
		})
		return
	}

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}
