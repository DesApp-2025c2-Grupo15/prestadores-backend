package main

import (
	"prestadores-api/internal/handler/afiliados"
	"prestadores-api/internal/handler/autorizaciones"
	"prestadores-api/internal/handler/login"
	"prestadores-api/internal/handler/recetas"
	"prestadores-api/internal/handler/reintegros"
	"prestadores-api/internal/handler/situaciones"
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

	// Repository y Service de Situaciones terapéuticas
	situacionRepo := repository.NewSituacionRepository()
	situacionService := service.NewSituacionService(situacionRepo, logger)

	// Handlers
	loginHandler := login.NewLoginHandler(logger)
	afiliadosHandler := afiliados.NewAfiliadoHandler(logger)
	historiaHandler := afiliados.NewHistoriaClinicaHandler(logger)
	autorizacionHandler := autorizaciones.NewAutorizacionHandler(autorizacionService, logger)
	recetaHandler := recetas.NewRecetaHandler(recetaService, logger)
	reintegroHandler := reintegros.NewReintegroHandler(reintegroService, logger)
	situacionHandler := situaciones.NewSituacionHandler(situacionService, logger)

	// Rutas /v1/prestadores
	v1 := r.Group("/v1/prestadores")
	{
		// Healthcheck
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Login
		v1.POST("/login", loginHandler.Login)

		// Afiliados
		afiliadosGroup := v1.Group("/afiliados")
		{
			afiliadosGroup.GET("", afiliadosHandler.GetAfiliados)
			afiliado := afiliadosGroup.Group("/:afiliadoId")
			{
				afiliado.GET("", afiliadosHandler.GetAfiliadoDetalle)
				afiliado.GET("/historia-clinica", historiaHandler.GetHistoriaClinica)
				// Situaciones terapéuticas

				afiliado.GET("/situaciones", situacionHandler.GetSituaciones)                               // ?scope=grupo
				afiliado.POST("/situaciones", situacionHandler.CreateSituacion)                             // alta
				afiliado.PATCH("/situaciones/:situacionId", situacionHandler.PatchSituacion)                // ej. fechaFin
				afiliado.PATCH("/situaciones/:situacionId/estado", situacionHandler.CambiarEstadoSituacion) // ALTA/BAJA/ACTIVA
				afiliado.DELETE("/situaciones/:situacionId", situacionHandler.DeleteSituacion)              // baja física (opcional)
			}
		}

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
