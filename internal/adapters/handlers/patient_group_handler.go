package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPatientGroupsByCodeOrOrgTitle godoc
// @Summary Получить группы пациентов по их коду или названию их организаций
// @Description Возвращает список экстренных приёмов, назначенных врачу на указанную дату, с пагинацией
// @Tags SMP
// @Accept json
// @Produce json
// @Param doc_id path uint true "ID врача"
// @Param date query string false "Дата в формате YYYY-MM-DD"
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(5)
// @Success 200 {array} models.PatientGroupShortResponse "Список приёмов"
// @Failure 400 {object} IncorrectFormatError "Некорректный запрос"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка"
// @Router /emergency/{doc_id} [get]
// GetPatientGroupsByCodeOrOrgTitle — поиск групп пациентов по коду или названию организации
func (h *Handler) GetPatientGroupsByCodeOrOrgTitle(c *gin.Context) {
	// 1. Парсим page
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

	// 2. Парсим perPage
	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, err := h.service.ParseUintString(perPageStr)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'perPage' must be an integer", false)
		return
	}
	if perPage == 0 || perPage > 100 {
		h.ErrorResponse(c, errors.New("perPage must be between 1 and 100"), http.StatusBadRequest, "parameter 'perPage' must be between 1 and 100", true)
		return
	}

	// 3. Получаем search
	search := c.Query("search")

	// 4. Вызываем usecase
	groups, appErr := h.usecase.GetPatientGroupsByCodeOrOrgTitle(search, int(page), int(perPage))
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 5. Возвращаем результат
	h.ResultResponse(c, "Success fetch patient groups", Object, groups)
}

// GetPatientGroupsByOrganization godoc
// @Summary Получить группы пациентов организации
// @Description Возвращает список групп пациентов конкретной организации с пагинацией
// @Tags PatientGroups
// @Accept json
// @Produce json
// @Param org_id path uint true "ID организации"
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(10)
// @Success 200 {object} models.PatientGroupShortResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /organizations/{org_id}/groups [get]
func (h *Handler) GetPatientGroupsByOrganization(c *gin.Context) {
	// 1. Парсим org_id из URL
	orgID, err := h.service.ParseUintString(c.Param("org_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'org_id' must be an uint", false)
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
		h.ErrorResponse(c, errors.New("perPage must be between 1 and 100"), http.StatusBadRequest, "parameter 'perPage' must be between 1 and 100", true)
		return
	}

	// 4. Вызываем usecase
	groups, appErr := h.usecase.GetPatientGroupsByOrganizationID(orgID, int(page), int(perPage))
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 5. Возвращаем результат
	h.ResultResponse(c, "Success fetch patient groups by organization", Object, groups)
}
