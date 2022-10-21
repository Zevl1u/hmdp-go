package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
	"time"
)

func RecordTime(c *gin.Context) {
	start := time.Now().UnixMilli()
	c.Next()
	end := time.Now().UnixMilli()
	fmt.Printf("Request %s use %d MilliSecs\n", c.Request.URL.Path, end-start)
}

func Session(key string) gin.HandlerFunc {
	// 创建基于cookie的存储引擎， 这里的"secret"可以随意设置，是一个加密密钥
	store := cookie.NewStore([]byte("secret"))

	// 配置存储引擎相关参数
	// 这里若要配置过期时间 必须配置路径 不然path会自动变成/user 暂时不知道为啥
	store.Options(sessions.Options{MaxAge: 3600, Path: "/"})
	return sessions.Sessions(key, store)
}

func LoginInterceptor(c *gin.Context) {
	session := sessions.Default(c)
	if user := session.Get("user"); user != nil {
		c.Next()
	} else {
		c.Abort()
		c.JSON(401, beans.Result{ErrMsg: "未登录，请先登录！"})
	}
}
