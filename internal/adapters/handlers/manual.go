package handlers

import "github.com/gin-gonic/gin"

// GetAllManuals godoc
// @Summary Получить все справочники
// @Description Возвращает полный список справочных записей (типов анализов, вакцин, организаций и т.д.)
// @Tags Manuals
// @Produce json
// @Success 200 {object} ResultResponse{data=[]models.ManualResponse} "Список справочников"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /manuals [get]
func (h *Handler) GetAllManuals(c *gin.Context) {
	ctx := c.Request.Context()

	manuals, appErr := h.usecase.GetAllManuals(ctx)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Manual entries retrieved successfully", Array, manuals)
}
