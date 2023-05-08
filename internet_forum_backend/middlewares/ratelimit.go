package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap) // 效率，总容量
	return func(c *gin.Context) {
		// 取不到令牌就返回响应
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit .... ")
			c.Abort() // 直接返回响应不在进行后续对context的操作
			return
		}
		// 取到令牌就放行
		c.Next()
	}
}
