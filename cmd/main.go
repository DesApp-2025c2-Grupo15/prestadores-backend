package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"prestadores-api/internal/handler"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Iniciando servidor prestadores-api",
		zap.String("port", "8080"),
		zap.String("version", "1.0.0"))

	r := gin.Default()

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

	logger.Info("Servidor iniciado exitosamente")
	r.Run(":8080")
}
