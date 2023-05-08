package controller

import "internet_forum/models"

// 接口文档需要用到的model
// 接口文档返回的数据格式是一致的，但是具体的Data类型不一致

type _ResponsePostList struct {
	Code    string                  `json:"code"`    // 状态码
	Message string                  `json:"message"` // 提示信息
	Data    []*models.ApiPostDetail `json:"data"`    // 数据
}

type _ResponsePostDetail struct {
	Code    string                `json:"code"`    // 状态码
	Message string                `json:"message"` // 提示信息
	Data    *models.ApiPostDetail `json:"data"`    // 数据
}

type _ResponseCreatePost struct {
	Code    string `json:"code"`    // 状态码
	Message string `json:"message"` // 提示信息
}

type _ResponseLogin struct {
	UserID   string `json:"user_id"`       // 用户ID
	Username string `json:"username"`      // 用户名
	AToken   string `json:"access_token"`  // atoken
	RToken   string `json:"refresh_token"` // rtoken
}
