package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllOrganizations godoc
// @Summary Получить все организации
// @Description Возвращает список организаций с пагинацией
// @Tags Organizations
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(5)
// @Success 200 {object} models.OrganizationShortResponse "Список организаций"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID пользователя"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /organizations [get]
// GetAllOrganizations — получение списка организаций для врача
// GetAllOrganizations — получение списка организаций для врача
func (h *Handler) GetAllOrganizations(c *gin.Context) {
	// 1. Получаем doctor_id из контекста
	doctorIDAny, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, nil, http.StatusUnauthorized, "Doctor ID not found in context", false)
		return
	}

	// 2.1. Парсим doctor_id
	doctorID, err := h.service.ParseUint(doctorIDAny)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Invalid doctor ID", false)
		return
	}

	// 3. Парсим page
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

	// 4. Парсим perPage
	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := h.service.ParseUintString(perPageStr)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'perPage' must be an integer", false)
		return
	}
	if perPage == 0 {
		h.ErrorResponse(c, errors.New("perPage must be greater than 0"), http.StatusBadRequest, "parameter 'perPage' must be greater than 0", true)
		return
	}

	// 5. Вызываем usecase
	organizations, appErr := h.usecase.GetAllOrganizations(doctorID, int(page), int(perPage))
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 6. Возвращаем результат
	h.ResultResponse(c, "Success get organizations", Object, organizations)
}
