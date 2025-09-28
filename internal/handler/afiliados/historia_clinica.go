package afiliados

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ===== Modelos =====

type NotaTurno struct {
	ID          int       `json:"id"`
	Fecha       time.Time `json:"fecha"`
	PrestadorID int       `json:"prestadorId"`
	Texto       string    `json:"texto"`
}

type Turno struct {
	ID           int         `json:"id"`
	Fecha        time.Time   `json:"fecha"`
	Especialidad string      `json:"especialidad"`
	Estado       string      `json:"estado"`
	Notas        []NotaTurno `json:"notas"`
}

type HistoriaClinica struct {
	AfiliadoID int     `json:"afiliadoId"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	Total      int     `json:"total"`
	Turnos     []Turno `json:"turnos"`
}

// ===== Handler =====

type HistoriaClinicaHandler struct {
	logger *zap.Logger
}

func NewHistoriaClinicaHandler(logger *zap.Logger) *HistoriaClinicaHandler {
	return &HistoriaClinicaHandler{
		logger: logger,
	}
}

// GetHistoriaClinica GET /v1/prestadores/afiliados/:id/historia-clinica - Query opcional: ?prestadorId=45 → filtra notas de ese prestador
func (h *HistoriaClinicaHandler) GetHistoriaClinica(c *gin.Context) {
	h.logger.Info("Obteniendo historia clínica",
		zap.String("endpoint", "/afiliados/:id/historia-clinica"),
		zap.String("method", "GET"))

	turnos := []Turno{
		{
			ID:           500,
			Fecha:        time.Date(2025, 9, 20, 10, 0, 0, 0, time.UTC),
			Especialidad: "Clínica",
			Estado:       "RESERVADO",
			Notas: []NotaTurno{
				{ID: 10, Fecha: time.Date(2025, 9, 20, 10, 30, 0, 0, time.UTC), PrestadorID: 45, Texto: "Control general"},
			},
		},
		{
			ID:           501,
			Fecha:        time.Date(2025, 9, 25, 15, 0, 0, 0, time.UTC),
			Especialidad: "Kinesiología",
			Estado:       "ATENDIDO",
			Notas: []NotaTurno{
				{ID: 10, Fecha: time.Date(2025, 9, 20, 10, 30, 0, 0, time.UTC), PrestadorID: 45, Texto: "Control general"},
				{ID: 12, Fecha: time.Date(2025, 9, 25, 15, 45, 0, 0, time.UTC), PrestadorID: 55, Texto: "Ejercicios domiciliarios"},
				{ID: 14, Fecha: time.Date(2025, 9, 25, 15, 50, 0, 0, time.UTC), PrestadorID: 55, Texto: "Se pudo notar un leve problema en la rodilla izquierda"},
			},
		},
	}

	// Filtro por prestadorId
	if pidStr := c.Query("prestadorId"); pidStr != "" {
		if pid, err := strconv.Atoi(pidStr); err == nil && pid > 0 {
			for i := range turnos {
				filtradas := make([]NotaTurno, 0, len(turnos[i].Notas))
				for _, n := range turnos[i].Notas {
					if n.PrestadorID == pid {
						filtradas = append(filtradas, n)
					}
				}
				turnos[i].Notas = filtradas
			}
		}
	}

	historia := HistoriaClinica{
		AfiliadoID: 1, // mock, en real debería salir de c.Param("id")
		Page:       0,
		Size:       20,
		Total:      len(turnos),
		Turnos:     turnos,
	}

	c.JSON(http.StatusOK, historia)
}
