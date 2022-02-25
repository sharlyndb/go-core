# go-core
### 介绍
go-core 是 go web 应用开发脚手架，从全局配置文件读取，zap日志组件始化，gorm数据库连接初始化，redis客户端初始化，http server启动等。最终实现简化流程、提高效率、统一规范。
### 安装
```
go get -u github.com/goworkeryyt/go-core
```
### 例子
默认的程序根目录下必须包含 resources 文件夹，且文件夹内必须有 active.yaml和四种不同环境的开发文件至少一种
配置文件参考 https://github.com/goworkeryyt/go-config 库的resources目录下的配置文件
```shell
├── resources(项目整合配置文件示例)
│   ├── active.yaml      配置指定要激活启用的配置文件
│   └── dev_config.yaml  开发环境配置文件
│   └── fat_config.yaml  功能验收测试环境配置文件
│   └── pro_config.yaml  生产环境配置文件
│   └── uat_config.yaml  用户验收测试环境配置文件
```
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-config"
	"github.com/goworkeryyt/go-config/env"
	"github.com/goworkeryyt/go-core/db"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-core/mqtt"
	"github.com/goworkeryyt/go-core/redis"
	"github.com/goworkeryyt/go-core/srun"
	"github.com/goworkeryyt/go-core/zap"
)

func main() {

	// 获取程序运行环境，默认会读取 resources/active.yaml 文件中配置的运行环境
	global.ENV = env.Active()

	// 获取全局配置,默认根据运行环境加载对应配置文件
	global.CONFIG = goconfig.GlobalConfig()

	// 初始化zap日志
	global.LOG = zap.Zap()

	// 初始化数据库连接
	global.DB = db.Gorm()

	// 初始化 redis 客户端
	global.REDIS = redis.Redis()

	// 初始化 mqtt
	global.MQTT = mqtt.DefaultMqtt("111111")

	// 获取配置文件原始内容,这样方便在程序中全局拿到自己定义的配置子项
	global.VP = global.CONFIG.Viper

	// 启动 http 服务
	r := gin.Default()
	// 健康监测
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
	// 启动服务
	srun.RunHttpServer(r)
}

```