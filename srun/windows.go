//go:build windows
//+build windows

/**
 * @Time: 2022/2/24 10:43
 * @Author: yt.yin
 */

package srun

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-core/global"
	"go.uber.org/zap"
)

// RunHttpServer Windows环境下启动服务
func RunHttpServer(r *gin.Engine) {
	address := fmt.Sprintf(":%d", global.CONFIG.Server.Addr)
	s := initServer(address, r)
	// 保证文本能够顺序输出
	time.Sleep(20 * time.Microsecond)
	global.LOG.Info("server run success on ", zap.String("address", address))
	err := s.ListenAndServe()
	if err != nil {
		global.LOG.Error(err.Error())
	}
}

// 初始化服务
func initServer(address string, router *gin.Engine)server {
	s := &http.Server{
		Addr:           address,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}
	return s
}




