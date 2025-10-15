package autorizaciones

import (
	"net/http"
	"prestadores-api/internal/model"
	"prestadores-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutorizacionHandler struct {
	service service.AutorizacionService
	logger  *zap.Logger
}

func NewAutorizacionHandler(service service.AutorizacionService, logger *zap.Logger) *AutorizacionHandler {
	return &AutorizacionHandler{
		service: service,
		logger:  logger,
	}
}

// Query params: estado?, q?, page?, size?, sort?
func (h *AutorizacionHandler) GetAutorizaciones(c *gin.Context) {
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

	h.logger.Info("Obteniendo lista de autorizaciones",
		zap.String("endpoint", "/solicitudes/autorizaciones"),
		zap.String("method", "GET"),
		zap.String("estado", estado),
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
	)

	response, err := h.service.GetAutorizaciones(estado, query, page, size, sort)
	if err != nil {
		h.logger.Error("Error al obtener autorizaciones", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener autorizaciones"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AutorizacionHandler) GetAutorizacionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	h.logger.Info("Obteniendo detalle de autorización",
		zap.String("endpoint", "/solicitudes/autorizaciones/:id"),
		zap.String("method", "GET"),
		zap.Int("id", id),
	)

	detalle, err := h.service.GetAutorizacionByID(id)
	if err != nil {
		h.logger.Error("Error al obtener autorización", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Autorización no encontrada"})
		return
	}

	c.JSON(http.StatusOK, detalle)
}

func (h *AutorizacionHandler) CreateAutorizacion(c *gin.Context) {
	var req model.CreateAutorizacionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para crear autorización", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Creando nueva autorización",
		zap.String("endpoint", "/solicitudes/autorizaciones"),
		zap.String("method", "POST"),
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.String("procedimiento", req.Procedimiento),
		zap.String("especialidad", req.Especialidad),
	)

	response, err := h.service.CreateAutorizacion(req)
	if err != nil {
		h.logger.Error("Error al crear autorización", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear autorización"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AutorizacionHandler) UpdateAutorizacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("ID inválido", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req model.UpdateAutorizacionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Request inválido para actualizar autorización", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request inválido", "details": err.Error()})
		return
	}

	h.logger.Info("Actualizando autorización",
		zap.String("endpoint", "/solicitudes/autorizaciones/:id"),
		zap.String("method", "PATCH"),
		zap.Int("id", id),
		zap.String("procedimiento", req.Procedimiento),
		zap.String("especialidad", req.Especialidad),
	)

	err = h.service.UpdateAutorizacion(id, req)
	if err != nil {
		h.logger.Error("Error al actualizar autorización", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Autorización no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Autorización actualizada exitosamente"})
}

func (h *AutorizacionHandler) CambiarEstadoAutorizacion(c *gin.Context) {
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

	h.logger.Info("Cambiando estado de autorización",
		zap.String("endpoint", "/solicitudes/autorizaciones/:id/estado"),
		zap.String("method", "PATCH"),
		zap.Int("id", id),
		zap.String("nuevoEstado", string(req.NuevoEstado)),
		zap.String("usuario", req.Usuario),
	)

	response, err := h.service.CambiarEstadoAutorizacion(id, req)
	if err != nil {
		h.logger.Error("Error al cambiar estado", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}