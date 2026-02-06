package main

import (
	"log"
	"todos_api/internal/config"
	"todos_api/internal/database"
	"todos_api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	var cfg *config.Config
	var err error

	cfg, err = config.Load() // create an instance of the config load

	if err != nil {
		log.Fatal("Unable to load config: %v", err)
	}

	var pool *pgxpool.Pool

	pool, err = database.Connect(cfg.DatabaseURL) // now create a pool form the config created
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	defer pool.Close()
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Welcome to GoLang To-Do REST API with Auth",
			"success":  true,
			"database": "connected",
		})

	})
	println("server starting")

	router.POST("/todos", handlers.CreateToDoHandler(pool))
	router.Run(":" + cfg.Port)

}
