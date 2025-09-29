package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginHandler struct {
	logger *zap.Logger
}

type LoginInfo struct {
	Username string `json:"username" binding:"required"`
}

func NewLoginHandler(logger *zap.Logger) *LoginHandler {
	return &LoginHandler{
		logger: logger,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	h.logger.Info("Ejecutando Login App",
		zap.String("endpoint", "/login"),
		zap.String("method", "POST"))

	var loginInfo LoginInfo

	if parseError := c.ShouldBindJSON(&loginInfo); parseError != nil {
		h.logger.Error("Error al parsear request", zap.Error(parseError))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de request inválido"})
		return
	}

	// username no puede ser vacío
	if loginInfo.Username == "" {
		h.logger.Warn("username vacío en login")
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'username' es obligatorio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login success",
		"username": loginInfo.Username,
	})
}
