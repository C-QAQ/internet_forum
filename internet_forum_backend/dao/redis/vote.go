package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票的权重
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	// 使用pipeline实现事务
	pipeline := client.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数(一天需要两百票才能将热度多持续一天)
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 判断投票限制
	// 取redis中帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 更新帖子分数
	// 先查询当前用户给当前帖子的投票记录
	oValue := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	if value == oValue {
		return ErrVoteRepeated
	}

	var op float64
	if value > oValue {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(value - oValue)

	pipeline := client.TxPipeline()
	// 更新帖子分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 记录用户为该帖子投过票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
