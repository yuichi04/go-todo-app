package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Todo struct {
	gorm.Model
	Content string `gorm:"not null"`
	Done    bool
	Until   time.Time
}

type Result struct {
	Content string `json:"content"`
}

func main() {
	connectionString := "host=localhost user=pquser password=pqpassword dbname=pqdatabase port=5432 sslmode=disable TimeZone=Asia/Tokyo"

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	e := echo.New()

	e.POST("/todos", func(c echo.Context) error {
		var todo Todo
		if err := c.Bind(&todo); err != nil {
			return err
		}
		if result := db.Create(&todo); result.Error != nil {
			return result.Error
		}
		return c.JSON(http.StatusCreated, todo)
	})

	e.GET("/todos", func(c echo.Context) error {
		var todos []Todo
		if result := db.Find(&todos); result.Error != nil {
			return result.Error
		}
		return c.JSON(http.StatusOK, todos)
	})

	e.GET("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		var todo Todo
		if result := db.First(&todo, id); result.Error != nil {
			return result.Error
		}
		var res Result
		res.Content = todo.Content
		return c.JSON(http.StatusOK, res)
	})

	e.PUT("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		var todo Todo
		if err := c.Bind(&todo); err != nil {
			return err
		}
		todo.ID = 0 // GORM が新しいレコードとして扱わないようにIDをリセット
		if result := db.Model(&Todo{}).Where("id = ?", id).Updates(todo); result.Error != nil {
			return result.Error
		}
		return c.JSON(http.StatusOK, todo)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		if result := db.Delete(&Todo{}, id); result.Error != nil {
			return result.Error
		}
		return c.NoContent((http.StatusNoContent))
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8989"))
}
