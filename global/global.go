/**
 * @Time: 2022/2/22 11:54
 * @Author: yt.yin
 */

package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"github.com/goworkeryyt/go-config"
	"github.com/goworkeryyt/go-config/env"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 单节点应用常用全局变量
var (

	// ENV 设置环境
	ENV  env.Environment

	// DB 数据库
	DB *gorm.DB

	// REDIS 默认客户端
	REDIS *redis.Client

	// MQTT 客户端
	MQTT *mqtt.Client

	// CONFIG 全局系统配置
	CONFIG *goconfig.Config

	// VP 通过 viper 读取的yaml配置文件
	VP *viper.Viper

	// LOG 全局日志
	LOG *zap.Logger

	// CSBEF casbin实施者
	CSBEF casbin.IEnforcer
)
