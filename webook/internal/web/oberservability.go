package web

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type ObservabilityHandler struct {
}

func (o *ObservabilityHandler) RegisterRouter(server *gin.Engine) {
	g := server.Group("test")
	g.GET("/prometheus", func(ctx *gin.Context) {
		sleep := rand.Int31n(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		ctx.String(http.StatusOK, "ok")
	})
}
