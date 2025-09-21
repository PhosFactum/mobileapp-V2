package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVersionProject возвращает версию проекта
// @Summary Получить версию проекта
// @Description Возвращает текущую версию приложения
// @Tags Utils
// @Produce json
// @Success 200 {object} map[string]string "Версия проекта"
// @Router /version [get]
func (h *Handler) GetVersionProject(c *gin.Context) {
	version := "1.2.3"
	// version := h.usecase.GetVersion()
	c.JSON(http.StatusOK, gin.H{"version": version})
}
