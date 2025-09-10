package main

import (
	afiliados "prestadores-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	afiliadosHandler := afiliados.NewAfiliadoHandler()

	v1 := r.Group("/v1/prestadores")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.GET("/afiliados", afiliadosHandler.GetAfiliados)
	}

	r.Run(":8080")
}
