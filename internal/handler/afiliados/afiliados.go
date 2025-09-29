package afiliados

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AfiliadoListItem struct {
	ID         int    `json:"id"`
	DNI        string `json:"dni"`
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	PlanMedico string `json:"planMedico"`
	Titular    bool   `json:"titular"`
}

// Detalle de afiliado
type AfiliadoDetalle struct {
	ID          int    `json:"id"`
	NroAfiliado string `json:"nroAfiliado"`
	DNI         string `json:"dni"`
	Nombre      string `json:"nombre"`
	Apellido    string `json:"apellido"`
	PlanMedico  string `json:"planMedico"`
	Titular     bool   `json:"titular"`
	Email       string `json:"email"`
	Telefono    string `json:"telefono"`
	Ciudad      string `json:"ciudad"`
	Provincia   string `json:"provincia"`
}

type AfiliadoHandler struct {
	logger *zap.Logger
}

func NewAfiliadoHandler(logger *zap.Logger) *AfiliadoHandler {
	return &AfiliadoHandler{logger: logger}
}

var afiliadosListMock = []AfiliadoListItem{
	{ID: 1, DNI: "43521489", Nombre: "María", Apellido: "Candia", PlanMedico: "Sancor Salud", Titular: true},
	{ID: 2, DNI: "53521489", Nombre: "Stella", Apellido: "Rodriguez", PlanMedico: "Galeno 210", Titular: true},
	{ID: 3, DNI: "40456015", Nombre: "Nicolas", Apellido: "Martin", PlanMedico: "Sancor Salud", Titular: false},
	{ID: 4, DNI: "12334555", Nombre: "Sofia", Apellido: "Lopez", PlanMedico: "Swiss Medical", Titular: false},
	{ID: 5, DNI: "11000189", Nombre: "Facundo", Apellido: "Gomez", PlanMedico: "Sancor Salud", Titular: true},
}

func detalleMockByID(id int) (AfiliadoDetalle, bool) {
	switch id {
	case 1:
		return AfiliadoDetalle{
			ID:          1,
			NroAfiliado: "15121231523",
			DNI:         "43521489",
			Nombre:      "María",
			Apellido:    "Candia",
			PlanMedico:  "Sancor Salud",
			Titular:     true,
			Email:       "juan.perez@email.com",
			Telefono:    "011-4567-8901",
			Ciudad:      "Buenos Aires",
			Provincia:   "Buenos Aires",
		}, true
	case 2:
		return AfiliadoDetalle{
			ID:          2,
			NroAfiliado: "15121231524",
			DNI:         "53521489",
			Nombre:      "Stella",
			Apellido:    "Candia",
			PlanMedico:  "Sancor Salud",
			Titular:     false,
			Email:       "maria.rodriguez@email.com",
			Telefono:    "0341-234-5678",
			Ciudad:      "Rosario",
			Provincia:   "Santa Fe",
		}, true
	default:
		return AfiliadoDetalle{}, false
	}
}

/* ===== Endpoints ===== */

// GET /v1/prestadores/afiliados  → devuelve un array simple para la tabla
func (h *AfiliadoHandler) GetAfiliados(c *gin.Context) {
	h.logger.Info("Obteniendo lista de afiliados",
		zap.String("endpoint", "/afiliados"),
		zap.String("method", "GET"))

	c.JSON(http.StatusOK, afiliadosListMock)
}

// GET /v1/prestadores/afiliados/:id  → devuelve detalle con grupo familiar
func (h *AfiliadoHandler) GetAfiliadoDetalle(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	h.logger.Info("Obteniendo detalle de afiliado",
		zap.String("endpoint", "/afiliados/:id"),
		zap.String("method", "GET"),
		zap.String("id", idStr),
	)

	detalle, ok := detalleMockByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Afiliado no encontrado"})
		return
	}

	c.JSON(http.StatusOK, detalle)
}
