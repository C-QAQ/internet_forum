package redis

import (
	"github.com/go-redis/redis"
	"internet_forum/models"
)

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

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	v1 := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}
	// 使用pipeline一次查询多个数据，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
