package main

import (
	"prestadores-api/internal/handler/afiliados"
	"prestadores-api/internal/handler/login"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	// Handlers
	loginHandler := login.NewLoginHandler(logger)
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
		v1.GET("/afiliados/:id", afiliadosHandler.GetAfiliadoDetalle)
		// Detalle de historia clínica (turnos + notas) del afiliado
		v1.GET("/afiliados/:id/historia-clinica", historiaHandler.GetHistoriaClinica)

		// Login
		v1.POST("/login", loginHandler.Login)
	}

	logger.Info("Servidor iniciado exitosamente")
	r.Run(":8080")
}
