package handler

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Suggest an ad for a client
// @Produce json
// @Success 200 {object} model.Ad
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 404 {object} ginerr.ErrorResp
// @Param clientId query string true "client_id"
// @Tags Ads
// @Router /ads [get]
func (h *Handler) getAd(c *gin.Context) {
	clientId, err := uuid.Parse(c.Query("client_id"))
	if err != nil {
		c.JSON(400, ginerr.Build("client_id must be uuid"))
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

	ad, err := h.adSvc.GetAd(client)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("no relevant ad found"))
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, ad)
}

// @Summary For testing: get all ad candidates, sorted in the order of priority
// @Produce json
// @Success 200 {object} []model.AdCandidate
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 404 {object} ginerr.ErrorResp
// @Param clientId query string true "client_id"
// @Tags Ads
// @Router /ads/candidates [get]
func (h *Handler) getAdCandidates(c *gin.Context) {
	clientId, err := uuid.Parse(c.Query("client_id"))
	if err != nil {
		c.JSON(400, ginerr.Build("client_id must be uuid"))
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

	candidates, err := h.adSvc.GetAdCandidates(client)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, candidates)
}

type adClickRequest struct {
	ClientId uuid.UUID `json:"client_id" binding:"required,uuid"`
}

// @Summary Notify that the ad was clicked
// @Produce json
// @Success 204
// @Failure 400 {object} ginerr.ErrorResp
// @Failure 404 {object} ginerr.ErrorResp
// @Failure 409 {object} ginerr.ErrorResp
// @Param adId path string true "adId"
// @Param request body adClickRequest true "request"
// @Tags Ads
// @Router /ads/{adId}/click [post]
func (h *Handler) clickAd(c *gin.Context) {
	var req adClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	campaign := c.MustGet("campaign").(model.Campaign)

	viewed, err := h.adSvc.IsAdViewed(req.ClientId, campaign.Id)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	if !viewed {
		c.JSON(409, ginerr.Build("ad was not viewed"))
		return
	}

	client, err := h.clientSvc.GetById(req.ClientId)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("client not found"))
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	if err := h.adSvc.ClickAd(client, campaign); err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.Status(204)
}
