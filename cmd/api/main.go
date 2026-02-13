package main

import (
	"log"
	"todos_api/internal/config"
	"todos_api/internal/database"
	"todos_api/internal/handlers"
	"todos_api/internal/middleware"

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

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", handlers.CreateToDoHandler(pool))
		protected.GET("", handlers.GetAllTodosHandler(pool))
		protected.GET("/:id", handlers.GetTodoByIDHandler(pool))
		protected.PUT("/:id", handlers.UpdateTodoHandler(pool))
		protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))
	}
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())
	router.Run(":" + cfg.Port)

}
