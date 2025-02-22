package handler

import (
	"backend/internal/model"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
	"strings"
)

// @Summary Upload image to campaign
// @Description Only .jpg and .png files up to 5 MB are allowed. This method won't fail if the campaign already has an image.
// @Produce json
// @Success 200 {object} model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param campaignId path string true "campaignId"
// @Param file formData file true "image"
// @Tags Images
// @Router /advertisers/{advertiserId}/campaigns/{campaignId}/image [put]
func (h *Handler) addCampaignImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	if !strings.HasSuffix(file.Filename, ".png") && !strings.HasSuffix(file.Filename, ".jpg") {
		c.JSON(400, ginerr.Build("image must be either png or jpg"))
		return
	}

	campaign := c.MustGet("campaign").(model.Campaign)
	campaign, err = h.imageSvc.AddCampaignImage(campaign, file)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, campaign)
}

// @Summary Delete image from campaign
// @Description Fails if the campaign does not have an image
// @Produce json
// @Success 200 {object} model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 404 {object} ginerr.ErrorResp
// @Param advertiserId path string true "advertiserId"
// @Param campaignId path string true "campaignId"
// @Tags Images
// @Router /advertisers/{advertiserId}/campaigns/{campaignId}/image [delete]
func (h *Handler) deleteCampaignImage(c *gin.Context) {
	campaign := c.MustGet("campaign").(model.Campaign)
	if campaign.ImagePath == "" {
		c.JSON(404, ginerr.Build("campaign does not have an image"))
		return
	}

	campaign, err := h.imageSvc.DeleteCampaignImage(campaign)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, campaign)
}
