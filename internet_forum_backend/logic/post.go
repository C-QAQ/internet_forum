package logic

import (
	"go.uber.org/zap"
	"internet_forum/dao/mysql"
	"internet_forum/dao/redis"
	"internet_forum/models"
	"internet_forum/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	p.ID = snowflake.GenID()
	if err = mysql.CreatePost(p); err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return err
}

// GetPostById 根据帖子id查询帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	var post *models.Post
	post, err = mysql.GetPostById(pid)
	if err != nil { // mysql层错误
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	var user *models.User
	user, err = mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID)",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	//根据社区id查询社区详细信息
	var community *models.CommunityDetail
	community, err = GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("GetCommunityDetailByID(post.CommunityID)",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	// 数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username, // 帖子作者名
		Post:            post,          // 帖子内容
		CommunityDetail: community,     // 帖子所在社区（板块）内容
	}
	return
}

// GetPostList 帖子分页查询
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		var user *models.User
		user, err = mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID)",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		var community *models.CommunityDetail
		community, err = GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID(post.CommunityID)",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	return
}

func GetPostListV2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// redis 查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data")
		return
	}
	zap.L().Debug("GetPostListV2", zap.Any("ids", ids))

	// 根据id去mysql数据库查询帖子详细信息
	// 返回的数据与给定的id的顺序一致
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs() failed", zap.Error(err))
		return
	}

	// 查询每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData() failed", zap.Error(err))
		return
	}

	// 拼接数据 将post的作者以及帖子所处社区信息绑定
	for idx, post := range posts {
		// 根据作者id查询作者信息
		var user *models.User
		user, err = mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID)",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		var community *models.CommunityDetail
		community, err = GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID(post.CommunityID)",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	return
}

// GetCommunityPostListV2 根据社区id返回帖子分页列表
func GetCommunityPostListV2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// redis 查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder() return 0 data")
		return
	}
	zap.L().Debug("GetCommunityListV2", zap.Any("ids", ids))

	// 根据id去mysql数据库查询帖子详细信息
	// 返回的数据与给定的id的顺序一致
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetCommunityPostListByIDs() failed", zap.Error(err))
		return
	}

	// 查询每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData() failed", zap.Error(err))
		return
	}

	// 拼接数据 将post的作者以及帖子所处社区信息绑定
	for idx, post := range posts {
		// 根据作者id查询作者信息
		var user *models.User
		user, err = mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID)",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		var community *models.CommunityDetail
		community, err = GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID(post.CommunityID)",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	return
}
