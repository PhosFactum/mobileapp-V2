package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// GetPersonalDataConsent возвращает PDF согласия на обработку персональных данных
// @Summary Получить PDF согласия на обработку персональных данных
// @Description Возвращает PDF-файл для ознакомления
// @Tags Consent
// @Security BearerAuth
// @Produce application/pdf
// @Success 200 {file} file "PDF документ"
// @Failure 500 {object} map[string]interface{} "Ошибка при чтении файла"
// @Router /consent/personal-data [get]
func (h *Handler) GetPersonalDataConsent(c *gin.Context) {
	filePath := "pkg/static/personal_data_consent.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		h.logger.Error("Failed to open PDF file", "path", filePath, "error", err)
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось открыть документ", false)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		h.logger.Error("Failed to get file info", "error", err)
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось получить информацию о файле", false)
		return
	}

	// ✅ Правильные заголовки для отображения в браузере
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\"personal_data_consent.pdf\"") // inline вместо attachment
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Cache-Control", "public, max-age=3600") // Кэширование на 1 час

	if _, err := io.Copy(c.Writer, file); err != nil {
		h.logger.Error("Failed to send PDF file", "error", err)
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось отправить документ", false)
		return
	}

	h.logger.Info("PDF displayed successfully", "file", filePath)
}

// GetMedicalExamConsent возвращает PDF согласия на медосмотр
// @Summary Получить PDF согласия на медицинский осмотр
// @Description Возвращает PDF-файл для ознакомления
// @Tags Consent
// @Security BearerAuth
// @Produce application/pdf
// @Success 200 {file} file "PDF документ"
// @Failure 500 {object} map[string]interface{} "Ошибка при чтении файла"
// @Router /consent/medical-exam [get]
func (h *Handler) GetMedicalExamConsent(c *gin.Context) {
	filePath := "pkg/static/medical_exam_consent.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		h.logger.Error("Failed to open PDF file", "path", filePath, "error", err)
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось открыть документ", false)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		h.logger.Error("Failed to get file info", "error", err)
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось получить информацию о файле", false)
		return
	}

	// ✅ Правильные заголовки для отображения в браузере
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\"medical_exam_consent.pdf\"") // inline вместо attachment
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Cache-Control", "public, max-age=3600") // Кэширование на 1 час

	if _, err := io.Copy(c.Writer, file); err != nil {
		h.logger.Error("Failed to send PDF file", "error", err)
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Не удалось отправить документ", false)
		return
	}

	h.logger.Info("PDF displayed successfully", "file", filePath)
}

// SaveSignature сохраняет подпись пациента
// @Summary Сохранить подпись пациента
// @Description Принимает изображение подписи и сохраняет его как подтверждение согласия
// @Tags Consent
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param recep_id path string true "ID пациента"
// @Param signature formData file true "Изображение подписи (PNG/JPG)"
// @Success 200 {object} entities.ConsentSignature "Подпись сохранена"
// @Failure 400 {object} IncorrectDataError "Неверный ID пациента или отсутствует файл"
// @Failure 500 {object} InternalServerError "Ошибка сервера"
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
// @Success 200 {string} signatureBase64.Base64 "Base64-кодированное изображение"
// @Failure 400 {object} IncorrectFormatError "Неверный ID пациента"
// @Failure 404 {object} NotFoundError "Подпись не найдена"
// @Failure 500 {object} InternalServerError "Ошибка сервера"
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
