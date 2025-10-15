package situaciones

import (
	"net/http"
	"prestadores-api/internal/model"
	"prestadores-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SituacionHandler struct {
	service service.SituacionService
	logger  *zap.Logger
}

func NewSituacionHandler(service service.SituacionService, logger *zap.Logger) *SituacionHandler {
	return &SituacionHandler{
		service: service,
		logger:  logger,
	}
}

// GET /v1/prestadores/afiliados/:afiliadoId/situaciones
// Soporta ?scope=grupo para devolver grupo familiar separado por integrante
func (h *SituacionHandler) GetSituaciones(c *gin.Context) {
	afiliadoIDStr := c.Param("afiliadoId")
	afiliadoID, err := strconv.Atoi(afiliadoIDStr)
	if err != nil || afiliadoID <= 0 {
		h.logger.Warn("afiliadoId inválido", zap.String("afiliadoId", afiliadoIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "afiliadoId inválido"})
		return
	}
	scope := c.DefaultQuery("scope", "") // "" | "grupo"

	h.logger.Info("Listando situaciones terapéuticas",
		zap.String("endpoint", "/afiliados/:afiliadoId/situaciones"),
		zap.String("method", "GET"),
		zap.Int("afiliadoId", afiliadoID),
		zap.String("scope", scope),
	)

	resp, err := h.service.GetSituaciones(afiliadoID, scope)
	if err != nil {
		h.logger.Error("Error al obtener situaciones", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener situaciones"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// POST /v1/prestadores/afiliados/:afiliadoId/situaciones
func (h *SituacionHandler) CreateSituacion(c *gin.Context) {
	afiliadoIDStr := c.Param("afiliadoId")
	afiliadoID, err := strconv.Atoi(afiliadoIDStr)
	if err != nil || afiliadoID <= 0 {
		h.logger.Warn("afiliadoId inválido", zap.String("afiliadoId", afiliadoIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "afiliadoId inválido"})
		return
	}

	var req model.CreateSituacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para crear situación", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}
	req.AfiliadoID = afiliadoID

	h.logger.Info("Creando nueva situación terapéutica",
		zap.String("endpoint", "/afiliados/:afiliadoId/situaciones"),
		zap.String("method", "POST"),
		zap.Int("afiliadoId", afiliadoID),
		zap.String("descripcion", req.Descripcion),
		zap.String("fechaInicio", req.FechaInicio),
	)

	resp, err := h.service.CreateSituacion(req)
	if err != nil {
		h.logger.Error("Error al crear situación", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear situación"})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// PATCH /v1/prestadores/afiliados/:afiliadoId/situaciones/:situacionId
// (Modificar datos, p.ej. fechaFin)
func (h *SituacionHandler) PatchSituacion(c *gin.Context) {
	afiliadoIDStr := c.Param("afiliadoId")
	situacionIDStr := c.Param("situacionId")
	_, err1 := strconv.Atoi(afiliadoIDStr)
	situacionID, err2 := strconv.Atoi(situacionIDStr)
	if err1 != nil || err2 != nil || situacionID <= 0 {
		h.logger.Warn("IDs inválidos", zap.String("afiliadoId", afiliadoIDStr), zap.String("situacionId", situacionIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	var req model.PatchSituacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para patch situación", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Patch situación terapéutica",
		zap.String("endpoint", "/afiliados/:afiliadoId/situaciones/:situacionId"),
		zap.String("method", "PATCH"),
		zap.Int("situacionId", situacionID),
	)

	if err := h.service.PatchSituacion(situacionID, req); err != nil {
		h.logger.Error("Error al patch situación", zap.Int("situacionId", situacionID), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Situación no encontrada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Situación actualizada exitosamente"})
}

// PATCH /v1/prestadores/afiliados/:afiliadoId/situaciones/:situacionId/estado
// (Dar de baja lógica, reactivar, etc.)
func (h *SituacionHandler) CambiarEstadoSituacion(c *gin.Context) {
	afiliadoIDStr := c.Param("afiliadoId")
	situacionIDStr := c.Param("situacionId")
	_, err1 := strconv.Atoi(afiliadoIDStr)
	situacionID, err2 := strconv.Atoi(situacionIDStr)
	if err1 != nil || err2 != nil || situacionID <= 0 {
		h.logger.Warn("IDs inválidos", zap.String("afiliadoId", afiliadoIDStr), zap.String("situacionId", situacionIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	var req model.CambioEstadoSituacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para cambiar estado de situación", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Cambiando estado de situación",
		zap.String("endpoint", "/afiliados/:afiliadoId/situaciones/:situacionId/estado"),
		zap.String("method", "PATCH"),
		zap.Int("situacionId", situacionID),
		zap.String("nuevoEstado", string(req.Estado)),
	)

	resp, err := h.service.CambiarEstadoSituacion(situacionID, req)
	if err != nil {
		h.logger.Error("Error al cambiar estado de situación", zap.Int("situacionId", situacionID), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DELETE /v1/prestadores/afiliados/:afiliadoId/situaciones/:situacionId
// (Baja física opcional)
func (h *SituacionHandler) DeleteSituacion(c *gin.Context) {
	afiliadoIDStr := c.Param("afiliadoId")
	situacionIDStr := c.Param("situacionId")
	_, err1 := strconv.Atoi(afiliadoIDStr)
	situacionID, err2 := strconv.Atoi(situacionIDStr)
	if err1 != nil || err2 != nil || situacionID <= 0 {
		h.logger.Warn("IDs inválidos", zap.String("afiliadoId", afiliadoIDStr), zap.String("situacionId", situacionIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	h.logger.Info("Eliminando situación terapéutica",
		zap.String("endpoint", "/afiliados/:afiliadoId/situaciones/:situacionId"),
		zap.String("method", "DELETE"),
		zap.Int("situacionId", situacionID),
	)

	if err := h.service.DeleteSituacion(situacionID); err != nil {
		h.logger.Error("Error al eliminar situación", zap.Int("situacionId", situacionID), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Situación no encontrada"})
		return
	}
	c.Status(http.StatusNoContent)
}
