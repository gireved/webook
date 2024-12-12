package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.Use(func(context *gin.Context) {
		println("这是第一个 Middleware")
	}, func(context *gin.Context) {
		println("这是第二个 Middleware")
	})
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})

	// 参数路由，路径参数
	server.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello, "+name)
	})

	// 查询参数
	// GET /order?id=123
	server.GET("/order", func(ctx *gin.Context) {
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "订单 ID 是 "+id)
	})

	server.GET("/views/*.html", func(ctx *gin.Context) {
		view := ctx.Param(".html")
		ctx.String(http.StatusOK, "view 是 "+view)
	})

	//server.GET("/star/*", func(ctx *gin.Context) {
	//	view := ctx.Param(".html")
	//	ctx.String(http.StatusOK, "view 是 "+view)
	//})
	//server.GET("/star/*/abc", func(ctx *gin.Context) {
	//	view := ctx.Param(".html")
	//	ctx.String(http.StatusOK, "view 是 "+view)
	//})

	server.POST("/login", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, login")
	})

	//go func() {
	//	server1 := gin.Default()
	//
	//	server1.Run(":8081")
	//}()

	// 如果你不传参数，那么实际上监听的是 8080 端口
	server.Run(":8080")
	// 这种写法是错的
	//  missing port in address
	//server.Run("8080")
}
