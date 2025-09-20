package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// GetAllPatients godoc
// @Summary Получить список всех пациентов
// @Description Возвращает список всех существующих пациентов с пагинацией
// @Tags Patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param count query int false "Количество записей" default(0)
// @Param filter query string false "Фильтр в формате field.operation.value"
// @Param order query string false "Сортировка в формате field.direction"
// @Success 200 {object} models.PatientFilterResponse "Список пациентов с пагинацией"
// @Failure 400 {object} IncorrectDataError "Некорректные данные"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка"
// @Router /patients [get]
func (h *Handler) GetAllPatients(c *gin.Context) {
	page, err := h.service.ParseIntString(c.DefaultQuery("page", "1"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'page' must be an integer", false)
		return
	}

	count, err := h.service.ParseIntString(c.DefaultQuery("count", "0"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'count' must be an integer", false)
		return
	}

	filter := c.Query("filter")
	order := c.Query("order")

	patients, appErr := h.usecase.GetAllPatients(page, count, filter, order)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, patients)
}

// CreatePatient godoc
// @Summary Создать нового пациента
// @Description Создает нового пациента
// @Tags Patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CreatePatientRequest true "Данные пациента"
// @Success 201 {object} entities.Patient "Созданный пациент"
// @Failure 400 {object} map[string]string "Неверный формат"
// @Failure 422 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /patients [post]
func (h *Handler) CreatePatient(c *gin.Context) {
	var input models.CreatePatientRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	patient, eerr := h.usecase.CreatePatient(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success patient create", Object, patient)
}
