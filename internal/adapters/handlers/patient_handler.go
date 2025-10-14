package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// GetAllPatients godoc
// @Summary Получить список всех пациентов
// @Description Возвращает список всех существующих пациентов
// @Description
// @Description Работает фильтрация, сортировка и пагинация
// @Tags Patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы\n(по умолчанию 1)"
// @Param count query int false "Количество записей на странице\n(по умолчанию 0 — без ограничения)"
// @Param filter query string false "Фильтр в формате field.operation.value.\nПримеры:\nfull_name.like.Иван - имя содержит 'Иван',\nbirth_date.eq.1988-07-14 - точная дата рождения"
// @Param order query string false "Сортировка в формате field.direction.\nПримеры:\nfull_name.asc - по алфавиту,\nid.desc - по убыванию ID пациента"
// @Success 200 {object} models.PatientsListResponse "Список пациентов"
// @Failure 400 {object} ResultError "Некорректные данные"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /patients [get]
func (h *Handler) GetPatientsByGroup(c *gin.Context) {
	group_id, err := h.service.ParseUintString(c.Param("group_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'group_id' must be an uint", false)
		return
	}
	ctx := c.Request.Context()
	patients, appErr := h.usecase.GetPatientsByGroup(ctx, group_id)

	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, patients)
}

// CreatePatient godoc
// @Summary Создать нового пациента
// @Description Создает нового пациента с персональными и контактными данными
// @Tags Patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param info body models.CreatePatientRequest true "Данные пациента"
// @Success 201 {object} entities.Patient "Созданный пациент"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /patients [post]
// handlers/patient_handler.go
// POST /patients - создание пациента
func (h *Handler) CreatePatient(c *gin.Context) {
	// 2. Биндим JSON
	var request models.CreatePatientRequest // ← обновлённое имя модели
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	// 3. Валидация через validator (если используется)
	if err := validate.Struct(request); err != nil {
		h.ErrorResponse(c, err, http.StatusUnprocessableEntity, "validation failed", true)
		return
	}

	// 4. Вызываем юзкейс с контекстом из Gin
	patient, appErr := h.usecase.CreatePatient(c.Request.Context(), request)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	// 6. Возвращаем результат
	h.ResultResponse(c, "Patient created successfully", Object, patient)
}
