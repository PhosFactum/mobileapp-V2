package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/gin-gonic/gin"
)

// GetAllOrganizations godoc
// @Summary Получить все организации
// @Description Возвращает список организаций с пагинацией
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(5)
// @Success 200 {object} models.OrganizationShortResponse "Список организаций"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID пользователя"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /organizations [get]
func (h *Handler) GetAllOrganizations(c *gin.Context) {
	doctorIDAny, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, nil, http.StatusUnauthorized, "Doctor ID not found in context", false)
		return
	}

	doctorID, err := services.ParseUint(doctorIDAny)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Invalid doctor ID", false)
		return
	}

	// Получаем номер страницы из query параметров
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.ErrorResponse(c, err, http.StatusBadRequest, "page must be a positive integer", false)
		return
	}

	// Получаем количество элементов на страницу из query параметров
	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		h.ErrorResponse(c, err, http.StatusBadRequest, "perPage must be a positive integer", false)
		return
	}

	// Вызываем usecase
	organizations, err := h.usecase.GetAllOrganizations(doctorID, page, perPage)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "error get organizations", false)
		return
	}
	h.ResultResponse(c, "Success get organizations", Object, organizations)
}
