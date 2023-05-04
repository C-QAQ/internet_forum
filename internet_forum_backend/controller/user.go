package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"internet_forum/dao/mysql"
	"internet_forum/logic"
	"internet_forum/models"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// zap 记录日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors翻译错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok { // 不是翻译错误，参数错误
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) // 翻译错误
		return
	}
	//2.业务处理
	if err := logic.Signup(p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) { // 用户已存在
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy) //	数据库错误 服务繁忙
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 获取请求参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam) // 参数错误
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) // 翻译错误
		return
	}
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseError(c, CodeServerBusy) // 服务错误
		return
	}
	ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID),
		"username":      user.Username,
		"access_token":  user.AToken,
		"refresh_token": user.RToken,
	})
}
