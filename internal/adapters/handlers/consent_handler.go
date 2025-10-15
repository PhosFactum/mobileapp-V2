package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// GetPersonalDataConsent godoc
// @Summary Получить PDF согласия на обработку персональных данных
// @Description Возвращает PDF-файл для ознакомления
// @Tags Consent
// @Security BearerAuth
// @Produce application/pdf
// @Success 200 "PDF документ"
// @Failure 500 {object} ResultError "Ошибка при чтении файла"
// @Router /consent/personal-data [get]
func (h *Handler) GetPersonalDataConsent(c *gin.Context) {
	filePath := "static/docs/consent_personal_data.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось открыть документ", false)
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\"consent_personal_data.pdf\"")

	if _, err := io.Copy(c.Writer, file); err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось отправить документ", false)
		return
	}
}

// GetMedicalExamConsent godoc
// @Summary Получить PDF согласия на медицинский осмотр
// @Description Возвращает PDF-файл для ознакомления
// @Tags Consent
// @Security BearerAuth
// @Produce application/pdf
// @Success 200 "PDF документ"
// @Failure 500 {object} ResultError "Ошибка при чтении файла"
// @Router /consent/medical-exam [get]
func (h *Handler) GetMedicalExamConsent(c *gin.Context) {
	filePath := "static/docs/consent_medical_exam.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось открыть документ", false)
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\"consent_medical_exam.pdf\"")

	if _, err := io.Copy(c.Writer, file); err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось отправить документ", false)
		return
	}
}

// SaveSignature godoc
// @Summary Сохранить подпись пациента
// @Description Принимает изображение подписи и сохраняет его как подтверждение согласия
// @Tags Consent
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param recep_id path integer true "ID приёма" minimum(1)
// @Param signature formData file true "Изображение подписи (PNG/JPG)"
// @Success 200 {object} map[string]interface{} "Подпись сохранена"
// @Failure 400 {object} ResultError "Неверный ID или отсутствует файл"
// @Failure 500 {object} ResultError "Ошибка сервера"
// @Router /consent/signature/{recep_id} [post]
func (h *Handler) SaveSignature(c *gin.Context) {
	patientID, err := h.service.ParseUintString(c.Param("recep_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	file, err := c.FormFile("signature")
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, errors.InternalServerError, false)
		return
	}
	defer openedFile.Close()

	signatureBytes, err := io.ReadAll(openedFile)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, errors.InternalServerError, false)
		return
	}

	// 👇 Используем ConsentUsecase через h.usecase
	if appErr := h.usecase.(interfaces.ConsentUsecase).SaveConsent(patientID, signatureBytes); appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Signature saved", Object, nil)
}

// GetSignature возвращает подпись пациента
// @Summary Получить подпись пациента
// @Description Возвращает ранее сохранённую подпись в формате Base64
// @Tags Consent
// @Produce json
// @Security BearerAuth
// @Param recep_id path string true "ID пациента"
// @Success 200 {object} map[string]interface{} "Подпись получена"
// @Success 200 {string} signatureBase64.Base64 "Base64-кодированное изображение"
// @Failure 400 {object} ResultError "Неверный ID пациента"
// @Failure 404 {object} ResultError "Подпись не найдена"
// @Failure 500 {object} ResultError "Ошибка сервера"
// @Router /consent/signature/{recep_id} [get]
func (h *Handler) GetSignature(c *gin.Context) {
	patientID, err := h.service.ParseUintString(c.Param("recep_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	// 👇 Используем ConsentUsecase через h.usecase
	signature, appErr := h.usecase.(interfaces.ConsentUsecase).GetSignature(patientID)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Signature fetched", Object, gin.H{"signatureBase64": signature})
}
