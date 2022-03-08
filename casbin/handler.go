/**
 * @Time: 2022/3/8 14:32
 * @Author: yt.yin
 */

package casbin

import (
	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-core/jwt"
	"net/http"
	"strings"
)

const (
	ADMI = "ADMI"
)

// CasbinHandler Casbin权限认证
func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.Contains(path, "swagger") {
			ctx.Next()
			return
		}
		if strings.Contains(path, "login") || strings.Contains(path, "health") || strings.Contains(path, "captcha") {
			ctx.Next()
			return
		}
		// 从ctx中获取claims
		claims, _ := jwt.GetClaims(ctx)
		user := claims.UserId
		permission := ctx.Request.URL.Path
		method := ctx.Request.Method

		if claims.UserType != ADMI {
			ok := CasbinServiceApp.PermissionVerify(user, permission, method)
			if !ok {
				ctx.JSON(http.StatusMethodNotAllowed, gin.H{
					"code":    -1,
					"message": "用户已经通过身份验证，但请求的接口:(" + permission + ")不在您的权限之内！",
				})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
		return
	}
}

