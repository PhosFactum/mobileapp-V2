package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// UpdateAnalysisOrder godoc
// @Summary Обновить направление на анализы
// @Description Обновляет список анализов в направлении
// @Tags AnalysisOrder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param info body models.UpdateAnalysisOrderRequest true "Данные направления"
// @Success 204 "Направление успешно обновлено"
// @Failure 400 {object} ResultError "Неверный формат запроса"
// @Failure 422 {object} ResultError "Ошибка валидации данных"
// @Failure 404 {object} ResultError "Направление не найдено"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /analysis-orders [patch]
func (h *Handler) UpdateAnalysisOrder(c *gin.Context) {
	// 1. Биндим JSON
	var request models.UpdateAnalysisOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	// 2. Вызываем юзкейс
	ctx := c.Request.Context()
	appErr := h.usecase.UpdateAnalysisOrder(ctx, &request)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 3. Успешный ответ
	c.Status(http.StatusNoContent)
}
