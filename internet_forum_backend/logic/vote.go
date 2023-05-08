package logic

import (
	"go.uber.org/zap"
	"internet_forum/dao/redis"
	"internet_forum/models"
	"strconv"
)

// 投票功能
// 1.用户投票的数据
// 2.

// 投一票就加432分 一天86400 / 200 -> 200张赞成票可以给一个帖子续一天热度

// VoteForPost 为帖子投票具体投票的业务逻辑
//
//	p.Direction=1时，有两种情况:
//		1. 之前没有投过票，现在投赞成票 --> 更新分数和投票纪录
//		2. 之前投反对票，现在改投赞成票 --> 更新分数和投票纪录
//	p.Direction=0时，有两种情况:
//		1. 之前投过赞成票，现在要取消投票 --> 更新分数和投票纪录
//		2. 之前投过反对票，现在要取消投票 --> 更新分数和投票纪录
//	p.Direction=-1时，有两种情况:
//		1. 之前没有投过票，现在投反对票 --> 更新分数和投票纪录
//		2. 之前投赞成票，现在改投反对票 --> 更新分数和投票纪录
//	投票限制：
//		每个帖子自发表之日起一个星期内允许用户投票，超过一个星期就不允许再投票。
//			1. 到期之后将redis中保存的赞成票数存储到mysql表中
//			2. 到期之后删除帖子 KeyPostVotedZSetPF
//	param userID *models.ParamVoteData
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("logic.VoteForPost()", zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
