package handler

import (
	"backend/internal/model"
	"backend/pkg/ginerr"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// @Summary Create campaign
// @Produce json
// @Success 200 {object} model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 409 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param request body model.CampaignCreateRequest true "request"
// @Tags Campaigns
// @Router /advertisers/{advertiserId}/campaigns [post]
func (h *Handler) createCampaign(c *gin.Context) {
	var req model.CampaignCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}
	if *req.StartDate > *req.EndDate {
		c.JSON(400, ginerr.Build("start date is after end date"))
		return
	}

	date := h.settingsSvc.Date()
	if *req.StartDate < date || *req.EndDate < date {
		c.JSON(409, ginerr.Build("either start or end date are in past"))
		return
	}

	adv := c.MustGet("advertiser").(model.Advertiser)

	start := time.Now()
	campaign, err := h.campaignSvc.Create(adv.Id, req)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	fmt.Println(time.Since(start))

	c.JSON(200, campaign)
}

// @Summary Get campaigns list
// @Produce json
// @Success 200 {object} []model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param size query int false "size"
// @Param page query int false "page"
// @Tags Campaigns
// @Router /advertisers/{advertiserId}/campaigns [get]
func (h *Handler) getCampaigns(c *gin.Context) {
	adv := c.MustGet("advertiser").(model.Advertiser)

	var req model.GetCampaignsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	if req.Size == 0 {
		req.Size = 100
	}
	if req.Page == 0 {
		req.Page = 1
	}

	campaigns, err := h.campaignSvc.GetPaginated(adv.Id, req.Size, req.Page)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, campaigns)
}

// @Summary Get campaign by id
// @Produce json
// @Success 200 {object} model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param campaignId path string true "campaignId"
// @Tags Campaigns
// @Router /advertisers/{advertiserId}/campaigns/{campaignId} [get]
func (h *Handler) getCampaignById(c *gin.Context) {
	c.JSON(200, c.MustGet("campaign"))
}

// @Summary Update campaign
// @Produce json
// @Success 200 {object} model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 409 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param campaignId path string true "campaignId"
// @Param request body model.CampaignCreateRequest true "request"
// @Tags Campaigns
// @Router /advertisers/{advertiserId}/campaigns/{campaignId} [put]
func (h *Handler) updateCampaign(c *gin.Context) {
	var req model.CampaignCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}
	if *req.StartDate > *req.EndDate {
		c.JSON(400, ginerr.Build("start date is after end date"))
		return
	}

	campaign := c.MustGet("campaign").(model.Campaign)
	date := h.settingsSvc.Date()

	if *req.StartDate != *campaign.StartDate && *req.StartDate < date {
		c.JSON(409, ginerr.Build("changed start date is in the past"))
		return
	}
	if *req.EndDate != *campaign.EndDate && *req.EndDate < date {
		c.JSON(409, ginerr.Build("changed end date is in the past"))
		return
	}

	if *campaign.StartDate <= date && (*campaign.StartDate != *req.StartDate ||
		*campaign.EndDate != *req.EndDate ||
		*campaign.ImpressionsLimit != *req.ImpressionsLimit ||
		*campaign.ClicksLimit != *req.ClicksLimit) {
		c.JSON(409, ginerr.Build("some of the updated fields can't be changed after campaign start"))
		return
	}

	err := h.campaignSvc.Update(&campaign, req)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, campaign)
}

// @Summary Delete campaign
// @Produce json
// @Success 204
// @Failure 400 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param campaignId path string true "campaignId"
// @Tags Campaigns
// @Router /advertisers/{advertiserId}/campaigns/{campaignId} [delete]
func (h *Handler) deleteCampaign(c *gin.Context) {
	campaign := c.MustGet("campaign").(model.Campaign)
	err := h.campaignSvc.Delete(campaign.Id)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.Status(204)
}
