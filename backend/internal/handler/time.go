package handler

import (
	"backend/internal/model"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
)

// @Summary Get current date
// @Produce json
// @Success 200 {object} model.CurrentDate
// @Tags Time
// @Router /time [get]
func (h *Handler) timeGet(c *gin.Context) {
	date := h.settingsSvc.Date()
	c.JSON(200, model.CurrentDate{CurrentDate: &date})
}

// @Summary Update current date
// @Produce json
// @Success 200 {object} model.CurrentDate
// @Failure 400 {object} ginerr.ErrorResp
// @Param request body model.CurrentDate true "request"
// @Tags Time
// @Router /time/advance [post]
func (h *Handler) timeAdvance(c *gin.Context) {
	var req model.CurrentDate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	err := h.settingsSvc.SetDate(*req.CurrentDate)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, req)
}
