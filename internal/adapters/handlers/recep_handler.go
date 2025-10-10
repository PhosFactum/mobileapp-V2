package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// CreateReception godoc
// @Summary Создать новый приём врача
// @Description Создаёт приём на основе шаблона и переданных данных
// @Tags Reception
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param info body models.CreateReceptionRequest true "Данные приёма"
// @Success 201 {object} entities.Reception "Созданный приём"
// @Failure 400 {object} errors.AppError "Неверный формат запроса или данные"
// @Failure 404 {object} errors.AppError "Шаблон или пациент не найден"
// @Failure 500 {object} errors.AppError "Внутренняя ошибка сервера"
// @Router /receptions [post]
func (h *Handler) CreateReception(c *gin.Context) {
	// 1. Биндим JSON
	var request models.CreateReceptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	// 2. Вызываем юзкейс с контекстом из Gin
	ctx := c.Request.Context()
	reception, appErr := h.usecase.CreateReception(ctx, &request)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 3. Возвращаем результат
	h.ResultResponse(c, "Reception created successfully", Object, reception)
}
