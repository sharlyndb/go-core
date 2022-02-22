/**
 * @Time: 2022/2/22 16:28
 * @Author: yt.yin
 */

package viper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
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
		fname := global.ENV.Value() + "_config"
		v.SetConfigName(fname)
		v.SetConfigType("yaml")
		v.AddConfigPath("./resources")
		log.Println("读取配置文件:", fname)
	} else {
		v.SetConfigFile(path[0])
		log.Println("读取指定配置文件:", path[0])
	}
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("严重错误的配置文件 : %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件内容发生改变:", e.Name)
		if err := v.Unmarshal(&global.CONFIGS); err != nil {
			log.Println("读取配置文件异常:", err)
		}
		global.CONFIGS.Viper = v
	})
	global.CONFIGS.Viper = v
	return v
}
