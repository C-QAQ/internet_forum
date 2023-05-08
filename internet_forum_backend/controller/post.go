package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"internet_forum/logic"
	"internet_forum/models"
	"strconv"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建新帖子，存入数据库并在redis中记录该帖子的分数和所处社区
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT_AToken"
// @Param obj body models.Post false "参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCreatePost
// @Router /api/v1/post [post]
func CreatePostHandler(c *gin.Context) {
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p)", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c 取出当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 存入数据库
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 通过post id获取post详情
// @Summary 通过post id获取post详情
// @Description 通过post id获取post内容以及所所在社区和作者名
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path int64 true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostDetail
// @Router /api/v1/post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 获取参数(从URL中获取帖子id)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil { // 帖子id格式错误
		zap.L().Error("get post detail whit invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetPostById(pid) // 调用logic层，获取帖子详情
	if err != nil {                     // logic层错误
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data) // 获取成功返回响应
}

// GetPostListHandler 帖子分页查询
// @Summary 概况
// @Description 描述
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param page path string false "页码"
// @Param size path string false "页面大小"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /api/v1/posts [post]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetPostListHandlerV2 获取帖子分页数据
// @Summary 获取帖子分页数据
// @Description 根据社区id（可以为空）、页码、数量返回分页数据
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /api/v2/posts [get]
func GetPostListHandlerV2(c *gin.Context) {
	op, err := strconv.ParseInt(c.Query("community_id"), 10, 64)
	if op > 0 && err == nil { // 如果query包含community_id则走另外一个handler
		GetCommunityPostListHandler(c)
		return
	}
	// 获取分页参数
	// 设置默认参数
	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
		CommunityID: 0,
	}
	// 绑定query参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandlerV2 with invalid params",
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("debug param info", zap.Any("param", p))
	// redis查询id列表
	data, err := logic.GetPostListV2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 获取帖子分页数据（含社区分类）
func GetCommunityPostListHandler(c *gin.Context) {
	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
		CommunityID: 0,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params",
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("debug param info", zap.Any("param", p))
	data, err := logic.GetCommunityPostListV2(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityListV2() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
