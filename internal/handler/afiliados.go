package afiliados

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Afiliado struct {
	ID        int    `json:"id"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	DNI       string `json:"dni"`
	Email     string `json:"email"`
	Telefono  string `json:"telefono"`
	Ciudad    string `json:"ciudad"`
	Provincia string `json:"provincia"`
}

type AfiliadoHandler struct{
	logger *zap.Logger
}

func NewAfiliadoHandler(logger *zap.Logger) *AfiliadoHandler {
	return &AfiliadoHandler{
		logger: logger,
	}
}

func (h *AfiliadoHandler) GetAfiliados(c *gin.Context) {
	h.logger.Info("Obteniendo lista de afiliados",
		zap.String("endpoint", "/afiliados"),
		zap.String("method", "GET"))

	afiliados := []Afiliado{
		{
			ID:        1,
			Nombre:    "Juan Carlos",
			Apellido:  "Perez",
			DNI:       "12345678",
			Email:     "juan.perez@email.com",
			Telefono:  "011-4567-8901",
			Ciudad:    "Buenos Aires",
			Provincia: "Buenos Aires",
		},
		{
			ID:        2,
			Nombre:    "Maria Elena",
			Apellido:  "Rodriguez",
			DNI:       "23456789",
			Email:     "maria.rodriguez@email.com",
			Telefono:  "0341-234-5678",
			Ciudad:    "Rosario",
			Provincia: "Santa Fe",
		},
		{
			ID:        3,
			Nombre:    "Carlos Alberto",
			Apellido:  "Gomez",
			DNI:       "34567890",
			Email:     "carlos.gomez@email.com",
			Telefono:  "0351-345-6789",
			Ciudad:    "Cordoba",
			Provincia: "Cordoba",
		},
		{
			ID:        4,
			Nombre:    "Ana Lucia",
			Apellido:  "Martinez",
			DNI:       "45678901",
			Email:     "ana.martinez@email.com",
			Telefono:  "0261-456-7890",
			Ciudad:    "Mendoza",
			Provincia: "Mendoza",
		},
		{
			ID:        5,
			Nombre:    "Roberto Miguel",
			Apellido:  "Silva",
			DNI:       "56789012",
			Email:     "roberto.silva@email.com",
			Telefono:  "0381-567-8901",
			Ciudad:    "San Miguel de Tucuman",
			Provincia: "Tucuman",
		},
	}

	h.logger.Info("Afiliados obtenidos exitosamente",
		zap.Int("total_afiliados", len(afiliados)),
		zap.String("status", "success"))

	c.JSON(http.StatusOK, gin.H{
		"data":    afiliados,
		"count":   len(afiliados),
		"message": "Afiliados obtenidos exitosamente",
	})
}
