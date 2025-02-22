package handler

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/ginerr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Get AI task status
// @Description Use short-polling with interval of 2 seconds to be notified when the task is completed
// @Produce json
// @Success 200 {object} model.AiTaskResponse
// @Failure 400 {object} ginerr.ErrorResp
// @Param taskId path string true "taskId"
// @Tags AI
// @Router /ai/tasks/{taskId} [get]
func (h *Handler) aiGetTask(c *gin.Context) {
	taskId, err := uuid.Parse(c.Param("taskId"))
	if err != nil {
		c.JSON(400, ginerr.Build("taskId must be uuid"))
		return
	}

	task, err := h.aiSvc.GetTask(taskId)
	if repo.IsNotFound(err) {
		c.JSON(404, ginerr.Build("task not found"))
		return
	}
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, task)
}

type aiSuggestTextRequest struct {
	AdTitle string `json:"ad_title" binding:"required"`
	Comment string `json:"comment"`
}

type aiSuggestTextResponse struct {
	TaskId uuid.UUID `json:"task_id"`
}

// @Summary Create a task to generate a list of suggestions
// @Description Create a task to suggest ad texts given by advertiser name and ad title. Returns task id to use in /ai/tasks/{taskId}
// @Produce json
// @Success 200 {object} aiSuggestTextResponse
// @Failure 400 {object} ginerr.ErrorResp
// @Param request body aiSuggestTextRequest true "request"
// @Param advertiserId path string true "advertiserId"
// @Tags AI
// @Router /ai/advertisers/{advertiserId}/suggestText [post]
func (h *Handler) aiSuggestText(c *gin.Context) {
	var req aiSuggestTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	adv := c.MustGet("advertiser").(model.Advertiser)

	taskId, err := h.aiSvc.SubmitSuggestText(adv.Name, req.AdTitle, req.Comment)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}
	c.JSON(200, aiSuggestTextResponse{TaskId: taskId})
}

// @Summary Get campaigns list with failed moderation
// @Produce json
// @Success 200 {object} []model.Campaign
// @Failure 400 {object} ginerr.ErrorResp
// @Param size query int false "size"
// @Param page query int false "page"
// @Tags Moderation
// @Router /ai/moderation/failed [get]
func (h *Handler) aiGetModerationFailed(c *gin.Context) {
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

	campaigns, err := h.campaignSvc.GetModerationFailed(req.Size, req.Page)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.JSON(200, campaigns)
}

type moderationStatus struct {
	Enabled *bool `json:"enabled" binding:"required"`
}

// @Summary Get moderation status
// @Produce json
// @Success 200 {object} moderationStatus
// @Tags Moderation
// @Router /ai/moderation/enabled [get]
func (h *Handler) aiModerationStatusGet(c *gin.Context) {
	enabled := h.settingsSvc.ModerationEnabled()
	c.JSON(200, moderationStatus{Enabled: &enabled})
}

// @Summary Enable/disable moderation (disabled by default)
// @Description When disabled, moderation_result field will be null
// @Produce json
// @Success 204
// @Failure 400 {object} ginerr.ErrorResp
// @Param request body moderationStatus true "request"
// @Tags Moderation
// @Router /ai/moderation/enabled [post]
func (h *Handler) aiModerationStatusUpdate(c *gin.Context) {
	var req moderationStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginerr.Build(err.Error()))
		return
	}

	err := h.settingsSvc.SetModerationEnabled(*req.Enabled)
	if err != nil {
		ginerr.Handle500(c, err)
		return
	}

	c.Status(204)
}
