package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

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
