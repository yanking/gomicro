// Package pprof provides pprof integration for the Gin HTTP framework.
// It registers the standard net/http/pprof handlers under a specified prefix.
package pprof

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

const (
	// defaultPrefix 默认pprof路径前缀
	defaultPrefix = "/debug/pprof"
)

// getPrefix 获取pprof路径前缀
func getPrefix(prefixOptions ...string) string {
	prefix := defaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

// Register 注册pprof路由到Gin引擎
func Register(r *gin.Engine, prefixOptions ...string) {
	// 获取路径前缀
	prefix := getPrefix(prefixOptions...)

	// 注册标准的 net/http/pprof 处理程序
	pprofGroup := r.Group(prefix)
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
}
