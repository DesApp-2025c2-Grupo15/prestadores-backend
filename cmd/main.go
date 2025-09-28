package main

import (
	"prestadores-api/internal/handler"
	"prestadores-api/internal/handler/afiliados"

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

	// Handlers
	loginHandler := handler.NewLoginHandler(logger)
	afiliadosHandler := afiliados.NewAfiliadoHandler(logger)
	historiaHandler := afiliados.NewHistoriaClinicaHandler(logger)

	// Rutas /v1/prestadores
	v1 := r.Group("/v1/prestadores")
	{
		// Healthcheck
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Afiliados
		v1.GET("/afiliados", afiliadosHandler.GetAfiliados)
		// Detalle de historia clínica (turnos + notas) del afiliado
		v1.GET("/afiliados/:id/historia-clinica", historiaHandler.GetHistoriaClinica)

		// Login
		v1.POST("/login", loginHandler.Login)
	}

	logger.Info("Servidor iniciado exitosamente")
	_ = r.Run(":8080")
}
