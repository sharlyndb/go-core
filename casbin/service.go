/**
 * @Time: 2022/3/8 14:37
 * @Author: yt.yin
 */

package casbin

import "github.com/goworkeryyt/go-core/global"

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

// AddPermissionForUserInDomain
/**
 *  @Description: 为用户或角色在域内添加权限
 *  @receiver casbinApi
 *  @param user 用户或角色
 *  @param domain 域
 *  @param permission url
 *  @param method 方法
 *  @return err
 */
func (s CasbinService) AddPermissionForUserInDomain(user, domain, permission, method string) (err error) {
	_, err = global.CSBEF.AddPermissionForUser(user, domain, permission, method)
	return
}

// PermissionVerify
/**
 *  @Description: 权限认证
 *  @receiver casbinApi
 *  @param user 用户或角色
 *  @param permission url
 *  @param method 方法
 *  @return ok 是否通过
 */
func (s CasbinService) PermissionVerify(user, permission, method string) (ok bool) {
	user = "user-" + user
	success, _ := global.CSBEF.Enforce(user, permission, method)
	if success {
		return true
	} else {
		return false
	}
}