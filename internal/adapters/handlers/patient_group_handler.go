package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPatientGroupsByOrganizationID godoc
// @Summary Получить группы пациентов по организации
// @Description Возвращает список групп пациентов, принадлежащих указанной организации, с возможностью поиска по коду группы и пагинацией
// @Tags PatientGroups
// @Accept json
// @Produce json
// @Param organization_id path uint true "ID организации"
// @Param search query string false "Поиск по коду группы (case-insensitive)"
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(10)
// @Success 200 {object} ResultResponse{data=[]models.PatientGroupShortResponse} "Список групп"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /patient-groups/by-organization/{organization_id} [get]
func (h *Handler) GetPatientGroupsByOrganizationID(c *gin.Context) {
	// 1. Получаем organization_id из пути
	orgIDStr := c.Param("organization_id")
	orgID, err := h.service.ParseUintString(orgIDStr)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'organization_id' must be a valid integer", false)
		return
	}
	if orgID == 0 {
		h.ErrorResponse(c, errors.New("organization_id must be greater than 0"), http.StatusBadRequest, "parameter 'organization_id' must be greater than 0", true)
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
	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, err := h.service.ParseUintString(perPageStr)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'perPage' must be an integer", false)
		return
	}
	if perPage == 0 || perPage > 100 {
		perPage = 10
	}

	// 4. Получаем search (опционально)
	search := c.Query("search")

	// 5. Вызываем usecase
	groups, appErr := h.usecase.GetPatientGroupsByOrganizationID(orgID, search, int(page), int(perPage))
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 6. Возвращаем результат
	h.ResultResponse(c, "Successfully fetched patient groups for organization", Object, groups)
}
