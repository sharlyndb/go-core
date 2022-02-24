//go:build !windows
// +build !windows

/**
 * @Time: 2022/2/24 10:49
 * @Author: yt.yin
 */

package srun

import (
	"fmt"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-core/global"
	"go.uber.org/zap"
)

// RunHttpServer Linux环境下启动服务
func RunHttpServer(r *gin.Engine) {
	address := fmt.Sprintf(":%d", global.CONFIG.Server.Addr)
	s := initUnixServer(address, r)
	// 保证文本顺序输出
	time.Sleep(20 * time.Microsecond)
	global.LOG.Info("server run success on ", zap.String("address", address))
	err := s.ListenAndServe()
	if err != nil {
		global.LOG.Error(err.Error())
	}
}

// 初始化服务
func initUnixServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
