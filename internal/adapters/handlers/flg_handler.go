package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateFLG godoc
// @Summary Создать запись ФЛГ
// @Description Добавляет новое флюорографическое исследование пациенту
// @Tags FLG
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param patient_id path uint true "ID пациента"
// @Param info body models.FLGCreateRequest true "Данные ФЛГ"
// @Success 201 {object} models.FLGResponse "Созданная запись ФЛГ"
// @Failure 400 {object} IncorrectFormatError "Некорректный запрос"
// @Failure 404 {object} NotFoundError "Пациент не найден"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка"
// @Router /patients/{patient_id}/flg [post]
func (h *Handler) CreateFLG(c *gin.Context) {
	patientID, err := services.ParseUint(c.Param("patient_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Некорректный ID пациента", true)
		return
	}

	var req models.FLGCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Некорректный формат запроса", true)
		return
	}

	req.PatientID = patientID
	resp, eerr := h.usecase.CreateFLG(&req)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "ФЛГ успешно создано", Object, resp)
}

// UpdateFLG godoc
// @Summary Обновить данные ФЛГ
// @Description Обновляет информацию о флюорографическом исследовании
// @Tags FLG
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param flg_id path uint true "ID записи ФЛГ"
// @Param info body models.FLGUpdateRequest true "Данные для обновления"
// @Success 200 {object} models.FLGResponse "Обновлённая запись ФЛГ"
// @Failure 400 {object} IncorrectFormatError "Некорректный запрос"
// @Failure 404 {object} NotFoundError "ФЛГ не найдено"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка"
// @Router /flg/{flg_id} [put]
func (h *Handler) UpdateFLG(c *gin.Context) {
	flgID, err := services.ParseUint(c.Param("flg_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Некорректный ID ФЛГ", true)
		return
	}

	var req models.FLGUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Некорректный формат запроса", true)
		return
	}

	resp, eerr := h.usecase.UpdateFLG(flgID, &req)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "ФЛГ успешно обновлено", Object, resp)
}
