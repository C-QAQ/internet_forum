package redis

import "internet_forum/models"

// GetPostIDsInOrder 从redis获取ids
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 根据用户请求中携带参数查询redis ids
	key := getRedisKey(KeyPostTimeZSet)

	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 确定查询的索引起止点

	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZREVRANGE 按照分数大小查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}
