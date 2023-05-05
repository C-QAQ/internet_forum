package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 定义请求参数结构体

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	// UserID 从请求中获取当前的用户
	// PostID 帖子id
	PostID string `json:"post_id" binding:"required"`
	// Direction 赞成票（1）还是反对票（-1）取消投票（0） 必须要有，只能是1 0 -1 其中一个
	Direction int8 `json:"direction,string" binding:"oneof=1 0 -1"`
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}
