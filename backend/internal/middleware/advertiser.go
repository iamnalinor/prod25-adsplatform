package middleware

import (
	"backend/internal/repo"
	"backend/internal/service"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdvertiserMiddleware struct {
	advertiserSvc *service.AdvertiserService
}

func NewAdvertiserMiddleware(advertiserSvc *service.AdvertiserService) *AdvertiserMiddleware {
	return &AdvertiserMiddleware{advertiserSvc}
}

func (m *AdvertiserMiddleware) Callback(c *gin.Context) {
	id, err := uuid.Parse(c.Param("advertiserId"))
	if err != nil {
		c.JSON(400, ginerr.Build("advertiserId must be uuid"))
		c.Abort()
		return
	}

	adv, err := m.advertiserSvc.GetById(id)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("advertiser not found"))
		c.Abort()
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		c.Abort()
		return
	}

	c.Set("advertiser", adv)
	c.Next()
}
