package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"hmdp/src/beans"
	"hmdp/src/utils"
	"hmdp/src/utils/db"
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

func RefreshTokenInterceptor(c *gin.Context) {
	ctx := c.Request.Context()
	auth := c.GetHeader(utils.AUTHORIZATION)
	// auth不为空才需要刷新，为空直接放行
	if auth != "" {
		m, err := db.RedisCli.HGetAll(ctx, utils.LOGIN_CODE_PREFIX+auth).Result()
		// 如果获取到的map不为空
		if len(m) != 0 { // 对于map的获取 不能err == redis.nil判断, 而要判断map是否为空map
			dto := beans.UserDTO{}
			utils.Map2Struct(m, &dto)
			c.Set("userDTO", dto)
			db.RedisCli.Expire(ctx, utils.LOGIN_CODE_PREFIX+auth, utils.LOGIN_USERDTO_TTL)
		} else if err != nil {
			panic(err)
		}
	}
	c.Next()

}

func LoginInterceptor(c *gin.Context) {
	v, exists := c.Get("userDTO")
	if exists {
		c.Next()
		fmt.Println(v)
	} else {
		c.Abort()
		c.JSON(401, beans.Result{ErrMsg: "未登录，请先登录！"})
	}
}
