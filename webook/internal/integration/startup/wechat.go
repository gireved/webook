package startup

import (
	"geektime-basic-go/webook/internal/service/oauth2/wechat"
	"geektime-basic-go/webook/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
