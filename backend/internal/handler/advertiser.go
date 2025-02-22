package handler

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
)

// @Summary Get advertiser by id
// @Produce json
// @Success 200 {object} model.Advertiser
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 404 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Tags Advertisers
// @Router /advertisers/{advertiserId} [get]
func (h *Handler) getAdvertiser(c *gin.Context) {
	adv := c.MustGet("advertiser").(model.Advertiser)
	c.JSON(200, adv)
}

// @Summary Upsert many advertisers at once
// @Produce json
// @Success 200 {object} []model.Advertiser
// @Failure 400 {object} ginerr.ErrorResp
// @Param request body []model.Advertiser true "request"
// @Tags Advertisers
// @Router /advertisers/bulk [post]
func (h *Handler) postAdvertisersBulk(c *gin.Context) {
	var req []model.Advertiser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	res, err := h.advertiserSvc.AddBulk(req)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, res)
}

// @Summary Add ML score for client-advertiser pair
// @Produce json
// @Success 200
// @Failure 400 {object} ginerr.ErrorResp
// @Param request body model.MlScore true "request"
// @Tags Advertisers
// @Router /ml-scores [post]
func (h *Handler) postMlScore(c *gin.Context) {
	var req model.MlScore
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	err := h.advertiserSvc.AddMlScore(req)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("client or advertiser not found"))
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.Status(200)
}
