package redis

// redis key

// redis key 注意使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix = "internet_forum:"
	// KeyPostTimeZSet ZSet:帖子及发帖时间
	KeyPostTimeZSet = "post:time"
	// KeyPostScoreZSet ZSet:帖子及投票的分数
	KeyPostScoreZSet = "post:score"
	// KeyPostVotedZSetPF ZSet:记录用户及投票的类型;参数是postID
	KeyPostVotedZSetPF = "post:voted:"
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
