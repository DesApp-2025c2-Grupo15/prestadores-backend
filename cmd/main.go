package main

import (
	"prestadores-api/internal/handler/afiliados"
	"prestadores-api/internal/handler/autorizaciones"
	"prestadores-api/internal/handler/login"
	"prestadores-api/internal/handler/recetas"
	"prestadores-api/internal/handler/reintegros"
	"prestadores-api/internal/repository"
	"prestadores-api/internal/service"
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

	// Repository y Service de autorizaciones
	autorizacionRepo := repository.NewAutorizacionRepository()
	autorizacionService := service.NewAutorizacionService(autorizacionRepo, logger)

	// Repository y Service de recetas
	recetaRepo := repository.NewRecetaRepository()
	recetaService := service.NewRecetaService(recetaRepo, logger)

	// Repository y Service de Reintegros
	reintegroRepo := repository.NewReintegroRepository()
	reintegroService := service.NewReintegroService(reintegroRepo, logger)

	// Handlers
	loginHandler := login.NewLoginHandler(logger)
	afiliadosHandler := afiliados.NewAfiliadoHandler(logger)
	historiaHandler := afiliados.NewHistoriaClinicaHandler(logger)
	autorizacionHandler := autorizaciones.NewAutorizacionHandler(autorizacionService, logger)
	recetaHandler := recetas.NewRecetaHandler(recetaService, logger)
	reintegroHandler := reintegros.NewReintegroHandler(reintegroService, logger)

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

		// Solicitudes
		solicitudes := v1.Group("/solicitudes")
		{
			// Autorizaciones
			autorizacionesGroup := solicitudes.Group("/autorizaciones")
			{
				autorizacionesGroup.GET("", autorizacionHandler.GetAutorizaciones)
				autorizacionesGroup.GET("/:id", autorizacionHandler.GetAutorizacionByID)
				autorizacionesGroup.POST("", autorizacionHandler.CreateAutorizacion)
				autorizacionesGroup.PATCH("/:id", autorizacionHandler.UpdateAutorizacion)
				autorizacionesGroup.PATCH("/:id/estado", autorizacionHandler.CambiarEstadoAutorizacion)
			}

			// Recetas
			recetasGroup := solicitudes.Group("/recetas")
			{
				recetasGroup.GET("", recetaHandler.GetRecetas)
				recetasGroup.GET("/:id", recetaHandler.GetRecetaByID)
				recetasGroup.POST("", recetaHandler.CreateReceta)
				recetasGroup.PUT("/:id", recetaHandler.UpdateReceta)
				recetasGroup.PATCH("/:id/estado", recetaHandler.CambiarEstadoReceta)
			}

			// Reintegros
			reintegrosGroup := solicitudes.Group("/reintegros")
			{
				reintegrosGroup.GET("", reintegroHandler.GetReintegros)
				reintegrosGroup.GET("/:id", reintegroHandler.GetReintegroByID)
				reintegrosGroup.POST("", reintegroHandler.CreateReintegro)
				reintegrosGroup.PUT("/:id", reintegroHandler.UpdateReintegro)
				reintegrosGroup.PATCH("/:id/estado", reintegroHandler.CambiarEstadoReintegro)
			}
		}
	}

	logger.Info("Servidor iniciado exitosamente")
	r.Run(":8080")
}
