package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) GetAllManuals(c *gin.Context) {
	ctx := c.Request.Context()

	manuals, appErr := h.usecase.GetAllManuals(ctx)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Manual entries retrieved successfully", Array, manuals)
}
