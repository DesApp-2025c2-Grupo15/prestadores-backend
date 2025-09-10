package main

import (
	afiliados "prestadores-api/internal/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Iniciando servidor prestadores-api", 
		zap.String("port", "8080"),
		zap.String("version", "1.0.0"))

	r := gin.Default()

	afiliadosHandler := afiliados.NewAfiliadoHandler(logger)

	v1 := r.Group("/v1/prestadores")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.GET("/afiliados", afiliadosHandler.GetAfiliados)
	}

	logger.Info("Servidor iniciado exitosamente")
	r.Run(":8080")
}
