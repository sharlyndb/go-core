/**
 * @Time: 2022/3/7 10:24
 * @Author: yt.yin
 */

package casbin

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/goworkeryyt/go-core/global"
)


// Casbin
/**
 *  @Description: 初始化Casbin执行者（与gorm结合）
 *  @param DB Gorm连接池
 *  @param casbinConfPath casbin配置文件地址
 *  @return Enforcer casbin执行者
 *  @return err 错误
 */
func Casbin(casbinModelPath ...string) *casbin.SyncedEnforcer{
	if global.DB == nil {
		log.Fatalln("未初始化数据库连接")
		return nil
	}
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		log.Fatalln("创建Casbin Gorm适配器错误："+err.Error())
		return nil
	}
	if len(casbinModelPath) > 0 {
		syncedEnforcer, err := casbin.NewSyncedEnforcer(casbinModelPath[0], adapter)
		if err != nil {
			log.Fatalln("初始化Casbin执行者错误："+err.Error())
			return nil
		}
		return syncedEnforcer
	}else{
		m, _ := model.NewModelFromString(`
			[request_definition]
			r = sub, obj, act
			
			[policy_definition]
			p = sub, obj, act
			
			[role_definition]
			g = _, _
			
			[policy_effect]
			e = some(where (p.eft == allow))
			
			[matchers]
			m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`)
		syncedEnforcer, err := casbin.NewSyncedEnforcer(m, adapter)
		if err != nil {
			log.Fatalln("初始化Casbin执行者错误："+err.Error())
			return nil
		}
		return syncedEnforcer
	}
}
