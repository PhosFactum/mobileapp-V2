package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllOrganizations godoc
// @Summary Получить организации врача
// @Description Возвращает список организаций, связанных с текущим врачом, с возможностью поиска по названию и пагинацией
// @Tags Organizations
// @Accept json
// @Produce json
// @Param search query string false "Поиск по названию организации"
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(5)
// @Success 200 {object} models.FilterResponse[[]models.OrganizationShortResponse] "Список организаций"
// @Failure 400 {object} errors.AppError "Неверный формат запроса"
// @Failure 401 {object} errors.AppError "Не авторизован"
// @Failure 500 {object} errors.AppError "Внутренняя ошибка сервера"
// @Router /organizations [get]
func (h *Handler) GetAllDoctorOrganizations(c *gin.Context) {
	// 1. Получаем doctor_id из контекста
	doctorIDAny, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, nil, http.StatusUnauthorized, "Doctor ID not found in context", false)
		return
	}

	doctorID, err := h.service.ParseUint(doctorIDAny)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Invalid doctor ID", false)
		return
	}

	// 2. Парсим page
	pageStr := c.DefaultQuery("page", "1")
	page, err := h.service.ParseUintString(pageStr)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'page' must be an integer", false)
		return
	}
	if page == 0 {
		h.ErrorResponse(c, errors.New("page must be greater than 0"), http.StatusBadRequest, "parameter 'page' must be greater than 0", true)
		return
	}

	// 3. Парсим perPage
	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := h.service.ParseUintString(perPageStr)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'perPage' must be an integer", false)
		return
	}
	if perPage == 0 || perPage > 100 {
		perPage = 5 // или верни ошибку, как в других местах
	}

	// 4. Получаем search (опционально)
	search := c.Query("search")

	// 5. Вызываем usecase
	organizations, appErr := h.usecase.GetAllDoctorOrganizations(doctorID, search, int(page), int(perPage))
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 6. Возвращаем результат
	h.ResultResponse(c, "Successfully fetched organizations", Object, organizations)
}
