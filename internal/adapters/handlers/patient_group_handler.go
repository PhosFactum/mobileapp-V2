package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPatientGroupsByCodeOrOrgTitle godoc
// @Summary Получить группы пациентов по их коду или названию их организаций
// @Description Возвращает список экстренных приёмов, назначенных врачу на указанную дату, с пагинацией
// @Tags Groups
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
func (h *Handler) GetPatientGroupsByCodeOrOrgTitle(c *gin.Context) {
	// Получаем дату из query параметров
	search := c.Query("search")

	// Получаем номер страницы
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	// Вызываем usecase
	receptions, err := h.usecase.GetPatientGroupsByCodeOrOrgTitle(search, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receptions)
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
	// Получаем org_id из URL
	orgIDParam := c.Param("org_id")
	orgID, err := strconv.ParseUint(orgIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Получаем параметры пагинации
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be positive integer"})
		return
	}

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PerPage must be between 1 and 100"})
		return
	}

	// Вызываем usecase
	groups, err := h.usecase.GetPatientGroupsByOrganizationID(uint(orgID), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}
