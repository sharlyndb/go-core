/**
 * @Time: 2022/3/7 11:58
 * @Author: yt.yin
 */

package jwt

import "github.com/golang-jwt/jwt/v4"

// CustomClaims 基础 claims
type CustomClaims struct {

	/** 单次登陆产生的tokenId */
	TokenId                 string         `json:"tokenId"           gorm:"column:token_id;index;comment:;type:varchar(36);"`

	/** 用户账号id */
	UserId                  string         `json:"userId"            gorm:"column:user_id;primary_key;comment:;type:varchar(36);"`

	/** 角色id兼容多个角色，用逗号分割 */
	AuthorityId             string         `json:"authorityId"       gorm:"column:authority_id;comment:角色ID;type:text;"`

	/** 用户类型 */
	UserType                string         `json:"userType"          gorm:"column:user_type;comment:用户类型;type:varchar(8)"`

	/** 用户名 */
	Username                string         `json:"username"          gorm:"column:username;comment:用户名;type:varchar(128);index;"`

	/** 用户昵称 */
	NickName                string         `json:"nickName"          gorm:"column:nick_name;comment:用户昵称;type:varchar(128);"`

	/** 电话 */
	Phone                   string         `json:"phone"             gorm:"column:phone;comment:电话;type:varchar(11)"`

	/** 用户所属商户的商户号 */
	MerchantNo              string         `json:"merchantNo"        gorm:"column:merchant_no;comment:商户号;type:varchar(32);"`

	/** 自定义扩展，可以为格式化后的json */
	Extend                  string         `json:"extend"            gorm:"column:extend;comment:自定义扩展，可以为格式化后的json;type:text;"`

	/** 有效时间 */
	BufferTime              int64          `json:"bufferTime"`

	/** 系统标准Claims */
	jwt.RegisteredClaims   `json:"-"        gorm:"-"`

}
