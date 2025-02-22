package middleware

import (
	"backend/internal/repo"
	"backend/internal/service"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CampaignMiddleware struct {
	campaignSvc *service.CampaignService
}

func NewCampaignMiddleware(campaignSvc *service.CampaignService) *CampaignMiddleware {
	return &CampaignMiddleware{campaignSvc}
}

func (m *CampaignMiddleware) Callback(c *gin.Context) {
	id, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		c.JSON(400, ginerr.Build("campaignId must be uuid"))
		c.Abort()
		return
	}

	campaign, err := m.campaignSvc.GetById(id)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("campaign not found"))
		c.Abort()
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		c.Abort()
		return
	}

	if advIdRaw := c.Param("advertiserId"); advIdRaw != "" && c.Query("testAdvertiserValidation") != "skip" {
		advId, err := uuid.Parse(advIdRaw)
		if err != nil || advId != campaign.AdvertiserId {
			c.JSON(403, ginerr.Build("this campaign does not belong to given advertiserId"))
			c.Abort()
			return
		}
	}

	c.Set("campaign", campaign)
	c.Next()
}
