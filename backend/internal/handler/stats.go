package handler

import (
	"backend/internal/model"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
)

// @Summary Get stats for campaign
// @Produce json
// @Success 200 {object} model.CampaignStats
// @Failure 404 {object} ginerr.ErrorResp
// @Param campaignId path string true "campaignId"
// @Tags Stats
// @Router /stats/campaigns/{campaignId} [get]
func (h *Handler) getStatsCampaign(c *gin.Context) {
	campaign := c.MustGet("campaign").(model.Campaign)
	stats, err := h.statsSvc.GetStatsCampaign(campaign)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, stats)
}

// @Summary Get stats for all campaigns of this advertiser
// @Produce json
// @Success 200 {object} model.CampaignStats
// @Failure 404 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Tags Stats
// @Router /stats/advertisers/{advertiserId}/campaigns [get]
func (h *Handler) getStatsAdvertiser(c *gin.Context) {
	adv := c.MustGet("advertiser").(model.Advertiser)
	stats, err := h.statsSvc.GetStatsAdvertiser(adv)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, stats)
}

// @Summary Get daily stats for campaign
// @Produce json
// @Success 200 {object} []model.CampaignStats
// @Failure 404 {object} ginerr.ErrorResp
// @Param campaignId path string true "campaignId"
// @Tags Stats
// @Router /stats/campaigns/{campaignId}/daily [get]
func (h *Handler) getStatsCampaignDaily(c *gin.Context) {
	campaign := c.MustGet("campaign").(model.Campaign)
	stats, err := h.statsSvc.GetStatsCampaignDaily(campaign)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, stats)
}

// @Summary Get daily stats for all campaigns of this advertiser
// @Produce json
// @Success 200 {object} []model.CampaignStats
// @Failure 404 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Tags Stats
// @Router /stats/advertisers/{advertiserId}/campaigns/daily [get]
func (h *Handler) getStatsAdvertiserDaily(c *gin.Context) {
	adv := c.MustGet("advertiser").(model.Advertiser)
	stats, err := h.statsSvc.GetStatsAdvertiserDaily(adv)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, stats)
}
