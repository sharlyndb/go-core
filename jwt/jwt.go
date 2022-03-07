/**
 * @Time: 2022/3/7 11:50
 * @Author: yt.yin
 */

package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/goworkeryyt/go-core/global"
	"gorm.io/gorm"
)

const (

	// DEF_SIGN_KEY 默认签名用的key
	DEF_SIGN_KEY= "82040620FEFAC4511FC65000ADAB0F77"
)

// 定义一些常量
var (
	TokenExpired     error = errors.New("Token 已经过期")
	TokenNotValidYet error = errors.New("Token 尚未激活")
	TokenMalformed   error = errors.New("Token 格式错误")
	TokenInvalid     error = errors.New("Token 无法解析")
)

// JWT jwt签名结构
type JWT struct {
	SigningKey []byte
}

// NewJWT 新建一个 jwt 实例
func NewJWT() *JWT {
	return &JWT{[]byte(GetSignKey())}
}

// GetSignKey 获取 signKey
func GetSignKey() string {
	key := global.CONFIG.JWT.SigningKey
	if key == "" {
		return DEF_SIGN_KEY
	}
	return key
}

// CreateToken 生成 token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 判断多点登录拦截是否开启
	if global.CONFIG.JWT.UseMultipoint {
		// 拦截
		if global.REDIS != nil {
			// 优先存入到 redis
			jsonData, _ := json.Marshal(claims)
			toJson :=string(jsonData)
			// 此处过期时间等于jwt过期时间
			timer := time.Duration(global.CONFIG.JWT.ExpiresTime) * time.Second
			err := global.REDIS.Set(context.Background(), claims.UserId, toJson, timer).Err()
			if err != nil {
				return "", err
			}
			return token.SignedString(j.SigningKey)
		}
		// 没有redis存入到 数据库
		err := global.DB.Save(&claims).Error
		if err != nil {
			return "", err
		} else {
			return token.SignedString(j.SigningKey)
		}
	}
	// 不拦截
	return token.SignedString(j.SigningKey)
}

// DeleteToken 强制删除Token记录，用途--用户账号被盗后，强制下线
func DeleteToken(userId string) (err error) {
	if global.REDIS != nil {
		err = global.REDIS.Del(context.Background(), userId).Err()
		return err
	}
	err = global.DB.Where("user_id = ?", userId).Delete(&CustomClaims{}).Error
	return err
}

// ResolveToken 解析token
func (j *JWT) ResolveToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// 判断多点登录拦截是否开启
			if global.CONFIG.JWT.UseMultipoint {
				clis := CustomClaims{}
				if global.REDIS != nil {
					jsonStr, err := global.REDIS.Get(context.Background(), claims.UserId).Result()
					if err == redis.Nil {
						return claims, nil
					}
					if jsonStr != "" {
						_ = json.Unmarshal([]byte(jsonStr), &clis)
						if clis.TokenId == "" {
							return claims, nil
						} else {
							if claims.TokenId == clis.TokenId {
								return claims, nil
							} else {
								return nil, errors.New("账号已在其他地方登录，您已被迫下线")
							}
						}
					}
					global.LOG.Error("从redis获取用户token异常：" + err.Error())
					return nil, errors.New("从redis获取用户token异常：" + err.Error())
				}
				err := global.DB.Where("user_id = ?", claims.UserId).First(&clis).Error
				if err != nil && err != gorm.ErrRecordNotFound {
					global.LOG.Error("从数据库获取用户token异常：" + err.Error())
					return nil, errors.New("从数据库获取用户token异常：" + err.Error())
				}
				if claims.TokenId == clis.TokenId {
					return claims, nil
				} else {
					return nil, errors.New("账号已在其他地方登录，您已被迫下线")
				}
			}
			return claims, nil
		}
	}
	return nil, TokenInvalid
}


// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Unix() + global.CONFIG.JWT.ExpiresTime
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// GetClaims 获取Claims
func GetClaims(c *gin.Context) (*CustomClaims, error) {
	if claims, exists := c.Get("claims"); !exists {
		global.LOG.Error("从Gin的Context中获取从jwt解析出来的用户claims失败, 请检查路由是否使用jwt中间件")
		return nil, errors.New("获取用户用户claims失败")
	} else {
		token := claims.(*CustomClaims)
		return token, nil
	}
}

// GetUserName 获取用户名称
func GetUserName(c *gin.Context) string {
	var username string
	if claims, ok := c.Get("claims"); ok {
		waitUse := claims.(*CustomClaims)
		username = waitUse.Username
	} else {
		username = ""
	}
	return username
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		global.LOG.Error("从Gin的Context中获取从jwt解析出来的用户UUID失败, 请检查路由是否使用jwt中间件!")
		return ""
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.AuthorityId
	}
}

// GetUserID 获取用户ID
func GetUserID(c *gin.Context) string {
	var userID string
	if claims, ok := c.Get("claims"); ok {
		waitUse := claims.(*CustomClaims)
		userID = waitUse.UserId
	} else {
		userID = ""
	}
	return userID
}

// GetTokenId 获取tokenId
func GetTokenId(c *gin.Context) string {
	var tokenId string
	if claims, ok := c.Get("claims"); ok {
		waitUse := claims.(*CustomClaims)
		tokenId = waitUse.TokenId
	} else {
		tokenId = ""
	}
	return tokenId
}

