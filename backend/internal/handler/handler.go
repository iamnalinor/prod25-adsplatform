package handler

import (
	"backend/config"
	"backend/docs"
	"backend/internal/middleware"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	adSvc         *service.AdService
	advertiserSvc *service.AdvertiserService
	aiSvc         *service.AiService
	apiSvc        *service.ApiService
	campaignSvc   *service.CampaignService
	clientSvc     *service.ClientService
	imageSvc      *service.ImageService
	settingsSvc   *service.SettingsService
	statsSvc      *service.StatsService
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		adSvc:         services.Ad,
		advertiserSvc: services.Advertiser,
		aiSvc:         services.Ai,
		apiSvc:        services.Api,
		campaignSvc:   services.Campaign,
		clientSvc:     services.Client,
		imageSvc:      services.Image,
		settingsSvc:   services.Settings,
		statsSvc:      services.Stats,
	}
}

func (h *Handler) GetRouter(env config.Environment) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 5 << 20 // limit for uploads: 5 MiB
	_ = router.SetTrustedProxies(nil)

	api := router.Group("")
	loggerMiddleware := middleware.NewRequestsLoggerMiddleware(h.apiSvc)
	if !env.RunningInCI {
		api.Use(loggerMiddleware.Callback)
	}

	advMiddleware := middleware.NewAdvertiserMiddleware(h.advertiserSvc)
	campaignMiddleware := middleware.NewCampaignMiddleware(h.campaignSvc)

	apiAdv := api.Group("")
	apiAdv.Use(advMiddleware.Callback)
	apiCampaign := api.Group("")
	apiCampaign.Use(campaignMiddleware.Callback)

	api.GET("/ping", h.ping)

	api.GET("/clients/:clientId", h.getClient)
	api.POST("/clients/bulk", h.postClientsBulk)

	apiAdv.GET("/advertisers/:advertiserId", h.getAdvertiser)
	api.POST("/advertisers/bulk", h.postAdvertisersBulk)
	api.POST("/ml-scores", h.postMlScore)

	apiAdv.POST("/advertisers/:advertiserId/campaigns", h.createCampaign)
	apiAdv.GET("/advertisers/:advertiserId/campaigns", h.getCampaigns)
	apiCampaign.GET("/advertisers/:advertiserId/campaigns/:campaignId", h.getCampaignById)
	apiCampaign.PUT("/advertisers/:advertiserId/campaigns/:campaignId", h.updateCampaign)
	apiCampaign.DELETE("/advertisers/:advertiserId/campaigns/:campaignId", h.deleteCampaign)

	api.GET("/ads", h.getAd)
	api.GET("/ads/candidates", h.getAdCandidates)
	apiCampaign.POST("/ads/:campaignId/click", h.clickAd)

	apiCampaign.GET("/stats/campaigns/:campaignId", h.getStatsCampaign)
	apiAdv.GET("/stats/advertisers/:advertiserId/campaigns", h.getStatsAdvertiser)
	apiCampaign.GET("/stats/campaigns/:campaignId/daily", h.getStatsCampaignDaily)
	apiAdv.GET("/stats/advertisers/:advertiserId/campaigns/daily", h.getStatsAdvertiserDaily)

	api.GET("/time", h.timeGet)
	api.POST("/time/advance", h.timeAdvance)

	apiCampaign.PUT("/advertisers/:advertiserId/campaigns/:campaignId/image", h.addCampaignImage)
	apiCampaign.DELETE("/advertisers/:advertiserId/campaigns/:campaignId/image", h.deleteCampaignImage)

	apiAdv.POST("/ai/advertisers/:advertiserId/suggestText", h.aiSuggestText)
	api.GET("/ai/tasks/:taskId", h.aiGetTask)
	api.GET("/ai/moderation/failed", h.aiGetModerationFailed)
	api.GET("/ai/moderation/enabled", h.aiModerationStatusGet)
	api.POST("/ai/moderation/enabled", h.aiModerationStatusUpdate)

	docs.SwaggerInfo.BasePath = "/"
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

// @Summary Ping the server
// @Produce json
// @Success 200 {object} map[string]string
// @Tags Ping
// @Router /ping [get]
func (h *Handler) ping(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
