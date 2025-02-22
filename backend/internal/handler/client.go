package handler

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Get client by id
// @Produce json
// @Success 200 {object} model.Client
// @Failure 400 {object} ginerr.ErrorResp
// @Param clientId path string true "clientId"
// @Tags Clients
// @Router /clients/{clientId} [get]
func (h *Handler) getClient(c *gin.Context) {
	clientId, err := uuid.Parse(c.Param("clientId"))
	if err != nil {
		c.JSON(400, ginerr.Build("clientId must be uuid"))
		return
	}

	client, err := h.clientSvc.GetById(clientId)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("client not found"))
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, client)
}

// @Summary Upsert many clients at once
// @Produce json
// @Success 200 {object} []model.Client
// @Failure 400 {object} ginerr.ErrorResp
// @Param request body []model.Client true "request"
// @Tags Clients
// @Router /clients/bulk [post]
func (h *Handler) postClientsBulk(c *gin.Context) {
	var req []model.Client
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	res, err := h.clientSvc.AddBulk(req)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, res)
}
