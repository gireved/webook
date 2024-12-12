package pkg

import (
	"geektime-basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapBody[T any](l logger.LoggerV1, fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		// 下半段业务逻辑哪里来
		// 我的业务逻辑有可能要操作 ctx
		// 你要读取HTTP HEADER
		res, err := fn(ctx, req)
		if err != nil {
			// 记录日志
			l.Error("处理业务逻辑出错", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
