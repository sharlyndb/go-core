/**
 * @Time: 2022/2/22 16:28
 * @Author: yt.yin
 */

package viper

import (
	"github.com/fsnotify/fsnotify"
	"github.com/goworkeryyt/go-config"
	"github.com/goworkeryyt/go-core/global"
	"github.com/spf13/viper"
	"log"
)

const (

	// ConfigSuffix 配置文件默认后缀
	ConfigSuffix = "_config"

	// ConfigType 配置文件类型
	ConfigType = "yaml"

	// ConfigPath 配置文件所在路径
	ConfigPath = "./resources"
)

// Viper 通过 viper 库读取 系统配置文件
func Viper(path ...string) *viper.Viper {
	v := viper.New()
	if len(path) == 0 {
		fname := global.ENV.Value() + ConfigSuffix
		v.SetConfigName(fname)
		v.SetConfigType(ConfigType)
		v.AddConfigPath(ConfigPath)
		log.Println("读取配置文件:", fname)
	} else {
		v.SetConfigFile(path[0])
		log.Println("读取指定配置文件:", path[0])
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置文件异常 : %s \n", err)
		return nil
	}
	global.CONFIG = &goconfig.Config{}
	if err := v.Unmarshal(global.CONFIG); err != nil {
		log.Fatalf("读取配置文件异常 : %s \n", err)
		return nil
	}
	global.CONFIG.Viper = v
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件内容发生改变:", e.Name)
		if err := v.Unmarshal(global.CONFIG); err != nil {
			log.Fatalf("读取配置文件异常 : %s \n", err)
			return
		}
		global.CONFIG.Viper = v
	})
	global.CONFIG.Viper = v
	return v
}

