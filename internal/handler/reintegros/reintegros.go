package reintegros

import (
	"net/http"
	"prestadores-api/internal/model"
	"prestadores-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReintegroHandler struct {
	service service.ReintegroService
	logger  *zap.Logger
}

func NewReintegroHandler(service service.ReintegroService, logger *zap.Logger) *ReintegroHandler {
	return &ReintegroHandler{
		service: service,
		logger:  logger,
	}
}

// Query params: estado?, q?, page?, size?, sort?
func (h *ReintegroHandler) GetReintegros(c *gin.Context) {
	estado := c.DefaultQuery("estado", "")
	query := c.DefaultQuery("q", "")
	pageStr := c.DefaultQuery("page", "0")
	sizeStr := c.DefaultQuery("size", "20")
	sort := c.DefaultQuery("sort", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 20
	}

	h.logger.Info("Obteniendo lista de reintegros",
		zap.String("endpoint", "/solicitudes/reintegros"),
		zap.String("method", "GET"),
		zap.String("estado", estado),
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
	)

	resp, err := h.service.GetReintegros(estado, query, page, size, sort)
	if err != nil {
		h.logger.Error("Error al obtener reintegros", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener reintegros"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ReintegroHandler) GetReintegroByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	h.logger.Info("Obteniendo detalle de reintegro",
		zap.String("endpoint", "/solicitudes/reintegros/:id"),
		zap.String("method", "GET"),
		zap.Int("id", id),
	)

	detalle, err := h.service.GetReintegroByID(id)
	if err != nil {
		h.logger.Error("Error al obtener reintegro", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Reintegro no encontrado"})
		return
	}
	c.JSON(http.StatusOK, detalle)
}

func (h *ReintegroHandler) CreateReintegro(c *gin.Context) {
	var req model.CreateReintegroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para crear reintegro", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Creando nuevo reintegro",
		zap.String("endpoint", "/solicitudes/reintegros"),
		zap.String("method", "POST"),
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.String("prestacion", req.Prestacion),
		zap.String("metodo", req.Metodo),
		zap.Float64("monto", req.Monto),
	)

	resp, err := h.service.CreateReintegro(req)
	if err != nil {
		h.logger.Error("Error al crear reintegro", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear reintegro"})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *ReintegroHandler) UpdateReintegro(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req model.UpdateReintegroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para actualizar reintegro", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Actualizando reintegro",
		zap.String("endpoint", "/solicitudes/reintegros/:id"),
		zap.String("method", "PUT"),
		zap.Int("id", id),
		zap.String("prestacion", req.Prestacion),
		zap.String("metodo", req.Metodo),
		zap.Float64("monto", req.Monto),
	)

	if err := h.service.UpdateReintegro(id, req); err != nil {
		h.logger.Error("Error al actualizar reintegro", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Reintegro no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reintegro actualizado exitosamente"})
}

func (h *ReintegroHandler) CambiarEstadoReintegro(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req model.CambioEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para cambiar estado", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Cambiando estado de reintegro",
		zap.String("endpoint", "/solicitudes/reintegros/:id/estado"),
		zap.String("method", "PATCH"),
		zap.Int("id", id),
		zap.String("nuevoEstado", string(req.NuevoEstado)),
		zap.String("usuario", req.Usuario),
	)

	resp, err := h.service.CambiarEstadoReintegro(id, req)
	if err != nil {
		h.logger.Error("Error al cambiar estado de reintegro", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
