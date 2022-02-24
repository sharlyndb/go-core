/**
 * @Time: 2022/2/24 14:59
 * @Author: yt.yin
 */

package consul

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
	"github.com/goworkeryyt/go-core/global"
)

var (
	// 注册中心
	cr  registry.Registry
	// 选择器
	sl  selector.Selector
)

// NewRegistry 初始化注册中心
func NewRegistry() registry.Registry {
	if cr != nil {
		return cr
	}
	cr = consul.NewRegistry(registry.Addrs(global.CONFIG.Consul.Addr))
	return cr
}

// NewRandomSelector 实例化节点选择器，策略是随机
func NewRandomSelector() selector.Selector {
	if sl != nil {
		return sl
	}
	// 实例化selector
	sl = selector.NewSelector(
		// 传入上面的consul
		selector.Registry(NewRegistry()),
		// 指定获取实例的算法
		selector.SetStrategy(selector.Random),
	)
	return sl
}

// NewRoundRobinSelector 实例化节点选择器，策略是轮询
func NewRoundRobinSelector() selector.Selector {
	if sl != nil {
		return sl
	}
	// 实例化selector
	sl = selector.NewSelector(
		// 传入上面的consul
		selector.Registry(NewRegistry()),
		// 指定获取实例的算法
		selector.SetStrategy(selector.RoundRobin),
	)
	return sl
}

