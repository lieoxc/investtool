// 更新 webserver/gin.go 以支持 API 模式
package webserver

import (
	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/investool/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go/extra"
)

func init() {
	// 替换 gin 默认的 validator，更加友好的错误信息
	binding.Validator = &goutils.GinStructValidator{}
	// causes the json binding Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	binding.EnableDecoderUseNumber = true

	// jsoniter 启动模糊模式来支持 PHP 传递过来的 JSON。容忍字符串和数字互转
	extra.RegisterFuzzyDecoders()
	// jsoniter 设置支持 private 的 field
	extra.SupportPrivateFields()
}

// NewGinEngine 根据参数创建 gin 的 router engine
// middlewares 需要使用到的中间件列表，默认不为 engine 添加任何中间件
func NewGinEngine() *gin.Engine {

	engine := gin.New()
	// ///a///b -> /a/b
	engine.RemoveExtraSlash = true

	// 添加默认中间件
	defaultMiddlewares := []gin.HandlerFunc{
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.CORS(),
	}

	// use middlewares
	for _, middleware := range defaultMiddlewares {
		engine.Use(middleware)
	}
	return engine
}
