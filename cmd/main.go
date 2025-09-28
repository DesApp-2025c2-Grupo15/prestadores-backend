package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"prestadores-api/internal/handler"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := gin.Default()

	// CORS: debe ir antes de las rutas
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vite
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // si no usás cookies, podés dejarlo en false
		MaxAge:           12 * time.Hour,
	}))

	loginHandler := handler.NewLoginHandler(logger)
	afiliadosHandler := handler.NewAfiliadoHandler(logger)

	v1 := r.Group("/v1/prestadores")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.GET("/afiliados", afiliadosHandler.GetAfiliados)
		v1.POST("/login", loginHandler.Login)
	}

	r.Run(":8080")
}
