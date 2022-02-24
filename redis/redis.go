/**
 * @Time: 2022/2/23 23:51
 * @Author: yt.yin
 */

package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/goworkeryyt/go-core/global"
	"go.uber.org/zap"
)

// Redis 初始z化redis客户端
func Redis() *redis.Client{
	redisCfg := global.CONFIG.Redis
	if redisCfg.Addr == "" {
		return nil
	}
	db := 0
	if redisCfg.DB >= 0 &&  redisCfg.DB <= 15{
		db = redisCfg.DB
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       db,
		MaxRetries:   redisCfg.MaxRetries,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
	})
	pong, err := client.Ping(context.TODO()).Result()
	if err != nil {
		global.LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
		return nil
	} else {
		global.LOG.Info("redis connect ping response:", zap.String("pong",pong))
		return client
	}
}
