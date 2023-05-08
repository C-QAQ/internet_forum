package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"internet_forum/logic"
	"internet_forum/models"
	"strconv"
)

// ---- 跟社区相关的 ----

// CommunityHandler 查询所有社区
// @Summary 查询所有社区
// @Description 查询到所有的社区 (community_id, community_name) 以列表的形式返回
// @Tags 社区
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Security ApiKeyAuth
// @Success 200
// @Router /api/v1/community [get]
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区 (community_id, community_name) 以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
// @Summary 概况
// @Description 描述
// @Tags 社区
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path string true "社区id"
// @Security ApiKeyAuth
// @Success 200
// @Router /api/v1/community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	var communityDetail *models.CommunityDetail
	communityDetail, err = logic.GetCommunityDetailByID(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetailByID(id) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, communityDetail)
}
