package logic

import (
	"internet_forum/dao/mysql"
	"internet_forum/models"
	"internet_forum/pkg/jwt"
	"internet_forum/pkg/snowflake"
)

func Signup(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 生成user_id（雪花算法）
	userID := snowflake.GenID()
	// 生成用户实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 插入用户到数据库
	err = mysql.InsertUser(user)
	return
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT token
	var aToken, rToken string
	aToken, rToken, err = jwt.GenToken(user.UserID, user.Username)
	user.AToken = aToken
	user.RToken = rToken
	return
}
