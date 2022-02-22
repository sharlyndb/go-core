/**
 * @Time: 2022/2/22 12:00
 * @Author: yt.yin
 */

package env

import (
    "fmt"
    "github.com/goworkeryyt/configs/profile"
    "github.com/spf13/viper"
    "log"
)

const (
    ConfigFile = "./resources/active.yaml"
)

// LoadDefaultActiveFile 读取默认配置文件 active.yaml
func LoadDefaultActiveFile() *profile.Profiles {
    v := viper.New()
    v.SetConfigFile(ConfigFile)
    err := v.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }
    var p profile.Profiles
    if err := v.Unmarshal(&p); err != nil {
        log.Println("Load active file err:", err)
        return nil
    }
    return &p
}
