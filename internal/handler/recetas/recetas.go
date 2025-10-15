package recetas

import (
	"net/http"
	"prestadores-api/internal/model"
	"prestadores-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecetaHandler struct {
	service service.RecetaService
	logger  *zap.Logger
}

func NewRecetaHandler(service service.RecetaService, logger *zap.Logger) *RecetaHandler {
	return &RecetaHandler{
		service: service,
		logger:  logger,
	}
}

// Query params: estado?, q?, page?, size?, sort?
func (h *RecetaHandler) GetRecetas(c *gin.Context) {
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

	h.logger.Info("Obteniendo lista de recetas",
		zap.String("endpoint", "/solicitudes/recetas"),
		zap.String("method", "GET"),
		zap.String("estado", estado),
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
	)

	response, err := h.service.GetRecetas(estado, query, page, size, sort)
	if err != nil {
		h.logger.Error("Error al obtener recetas", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener recetas"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *RecetaHandler) GetRecetaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	h.logger.Info("Obteniendo detalle de receta",
		zap.String("endpoint", "/solicitudes/recetas/:id"),
		zap.String("method", "GET"),
		zap.Int("id", id),
	)

	detalle, err := h.service.GetRecetaByID(id)
	if err != nil {
		h.logger.Error("Error al obtener receta", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Receta no encontrada"})
		return
	}

	c.JSON(http.StatusOK, detalle)
}

func (h *RecetaHandler) CreateReceta(c *gin.Context) {
	var req model.CreateRecetaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para crear receta", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Creando nueva receta",
		zap.String("endpoint", "/solicitudes/recetas"),
		zap.String("method", "POST"),
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.String("medicamento", req.Medicamento),
		zap.String("dosis", req.Dosis),
	)

	response, err := h.service.CreateReceta(req)
	if err != nil {
		h.logger.Error("Error al crear receta", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear receta"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *RecetaHandler) UpdateReceta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req model.UpdateRecetaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para actualizar receta", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Actualizando receta",
		zap.String("endpoint", "/solicitudes/recetas/:id"),
		zap.String("method", "PUT"),
		zap.Int("id", id),
		zap.String("medicamento", req.Medicamento),
		zap.String("dosis", req.Dosis),
	)

	err = h.service.UpdateReceta(id, req)
	if err != nil {
		h.logger.Error("Error al actualizar receta", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Receta no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Receta actualizada exitosamente"})
}

func (h *RecetaHandler) CambiarEstadoReceta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req model.CambioEstadoRecetaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para cambiar estado", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Cambiando estado de receta",
		zap.String("endpoint", "/solicitudes/recetas/:id/estado"),
		zap.String("method", "PATCH"),
		zap.Int("id", id),
		zap.String("nuevoEstado", string(req.NuevoEstado)),
		zap.String("usuario", req.Usuario),
	)

	response, err := h.service.CambiarEstadoReceta(id, req)
	if err != nil {
		h.logger.Error("Error al cambiar estado", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
