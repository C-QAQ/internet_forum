package logic

import (
	"internet_forum/dao/mysql"
	"internet_forum/models"
)

// GetCommunityList 查找到所有的community 并返回
func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
