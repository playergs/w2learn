package middleware

import (
	"strings"
	"w2learn/internal/utils"
	"w2learn/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func JWTAuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 检查 Authorization 头
		if authHeader == "" {
			response.Error(c, "Authorization header is empty")
			c.Abort()
			return
		}

		// 获取 token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		jwt, err := utils.ParseJWT(tokenStr)

		if err != nil {
			response.Error(c, "Invalid token")
			c.Abort()
			return
		}
		// 检查过期时间
		if utils.IsTokenExpired(jwt) {
			response.Error(c, "Token expired")
			c.Abort()
			return
		}

		// 检查 token 是否在 Redis 黑名单 中
		n, err := rdb.Exists(c.Request.Context(), jwt.ID).Result()

		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		if n > 0 {
			response.Error(c, "token is on the blacklist")
			c.Abort()
			return
		}

		c.Set("id", jwt.ID)

		c.Next()
	}
}
