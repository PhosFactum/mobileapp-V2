package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// CreateVaccine godoc
// @Summary Создать запись о вакцинации
// @Description Регистрирует факт вакцинации пациента с привязкой к справочникам
// @Tags Vaccines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateVaccineRequest true "Данные вакцинации"
// @Success 200 {object} handlers.ResultResponse{data=models.VaccineAllResponse} "Запись успешно создана"
// @Failure 400 {object} ResultError "Неверный формат запроса"
// @Failure 422 {object} ResultError "Ошибка валидации данных"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /vaccines [post]
func (h *Handler) CreateVaccine(c *gin.Context) {
	var req models.CreateVaccineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	vaccine, appErr := h.usecase.CreateVaccine(c.Request.Context(), &req)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Vaccine created successfully", Object, vaccine)
}

// CreateVaccineRefusal godoc
// @Summary Зафиксировать отказ от вакцинации
// @Description Регистрирует отказ пациента от вакцинации по определённому типу
// @Tags Vaccines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateVaccineRefusalRequest true "Данные отказа"
// @Success 200 {object} handlers.ResultResponse{data=models.VaccineAllResponse} "Отказ успешно зафиксирован"
// @Failure 400 {object} ResultError "Неверный формат запроса"
// @Failure 422 {object} ResultError "Ошибка валидации данных"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /vaccines/refusals [post]
func (h *Handler) CreateVaccineRefusal(c *gin.Context) {
	var req models.CreateVaccineRefusalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	refusal, appErr := h.usecase.CreateVaccineRefusal(c.Request.Context(), &req)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Vaccine refusal created successfully", Object, refusal)
}

// CreateVaccineWithdrawal godoc
// @Summary Зафиксировать отзыв вакцинации
// @Description Регистрирует отзыв вакцинации (например, при медотводе)
// @Tags Vaccines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateVaccineWithdrawalRequest true "Данные отзыва"
// @Success 200 {object} handlers.ResultResponse{data=models.VaccineAllResponse} "Отзыв успешно зафиксирован"
// @Failure 400 {object} ResultError "Неверный формат запроса"
// @Failure 422 {object} ResultError "Ошибка валидации данных"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /vaccines/withdrawals [post]
func (h *Handler) CreateVaccineWithdrawal(c *gin.Context) {
	var req models.CreateVaccineWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	withdrawal, appErr := h.usecase.CreateVaccineWithdrawal(c.Request.Context(), &req)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Vaccine withdrawal created successfully", Object, withdrawal)
}

// CreateTitr godoc
// @Summary Создать запись о титре антител
// @Description Регистрирует результат анализа на титр антител после вакцинации
// @Tags Vaccines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateTitrRequest true "Данные титра"
// @Success 200 {object} handlers.ResultResponse{data=models.VaccineAllResponse} "Титр успешно сохранён"
// @Failure 400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure 422 {object} handlers.ErrorResponse "Ошибка валидации данных"
// @Failure 500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router /vaccines/titrs [post]
func (h *Handler) CreateTitr(c *gin.Context) {
	var req models.CreateTitrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid request body", true)
		return
	}

	titr, appErr := h.usecase.CreateTitr(c.Request.Context(), &req)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Titr created successfully", Object, titr)
}
