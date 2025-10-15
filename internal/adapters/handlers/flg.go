package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// CreateFlgWithPhoto godoc
// @Summary Создать флюорографию с фото
// @Description Загружает фото и создаёт запись
// @Tags Flg
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param patient_id formData uint true "ID пациента"
// @Param organization formData string true "Организация"
// @Param number formData string true "Номер"
// @Param result formData string true "Результат"
// @Param date formData string true "Дата (YYYY-MM-DD)"
// @Param file formData file true "Фото (JPEG/PNG, до 10 МБ)"
// @Success 201 {object} models.FlgResponse
// @Failure 400 {object} ResultError
// @Failure 404 {object} ResultError
// @Failure 422 {object} ResultError
// @Failure 500 {object} ResultError
// @Router /flg [post]
func (h *Handler) CreateFlgWithPhoto(c *gin.Context) {
	// 1. Получение файла
	file, err := c.FormFile("file")
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "file is required", true)
		return
	}

	// 2. Парсинг остальных полей
	patientIDStr := c.PostForm("patient_id")
	organization := c.PostForm("organization")
	number := c.PostForm("number")
	result := c.PostForm("result")
	date := c.PostForm("date")

	// 3. Валидация и конвертация
	patientID, err := strconv.ParseUint(patientIDStr, 10, 64)
	if err != nil || patientID == 0 {
		h.ErrorResponse(c, nil, http.StatusUnprocessableEntity, "invalid patient_id", true)
		return
	}

	// 4. Валидация Content-Type
	contentType := file.Header.Get("Content-Type")
	if !isValidImageContentType(contentType) {
		h.ErrorResponse(c, nil, http.StatusUnprocessableEntity, "invalid content type: only JPEG/PNG allowed", true)
		return
	}

	// 5. Чтение файла в память
	src, err := file.Open()
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to open file", false)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to read file", false)
		return
	}

	// 6. Ограничение размера (дублируем валидацию для безопасности)
	if len(data) > 10<<20 {
		h.ErrorResponse(c, nil, http.StatusRequestEntityTooLarge, "file too large (max 10 MB)", true)
		return
	}

	// 7. Формирование запроса
	req := &models.CreateFlgRequest{
		PatientID:    uint(patientID),
		Organization: organization,
		Number:       number,
		Result:       result,
		Date:         date,
		FileData:     data,
		ContentType:  contentType,
	}

	// 8. Вызов юзкейса
	resp, appErr := h.usecase.CreateFlgWithPhoto(c.Request.Context(), req)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Вспомогательная функция (дублируем для хендлера)
func isValidImageContentType(ct string) bool {
	ct = strings.ToLower(ct)
	return ct == "image/jpeg" || ct == "image/jpg" || ct == "image/png"
}
