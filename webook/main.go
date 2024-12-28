package main

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	initViperV1()
	initLogger()
	initPrometheus()
	app := InitWebServer()
	// Consumer 类似于Web 或者GRPC，是一个顶级入口
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.server
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})
	server.Run(":8080")
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// 当前工作目录的 config 子目录
	viper.AddConfigPath("config")
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperWatch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
	//viper.Set("db.dsn", "localhost:3306")
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println(viper.GetString("test.key"))
	})
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperV1() {
	cfile := pflag.String("dev",
		"config/dev.yaml", "配置文件路径")
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
	//viper.Set("db.dsn", "localhost:3306")
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperV2() {
	cfg := `
test:
  key: value1

redis:
  addr: "localhost:6379"

db:
  dsn: "root:root@tcp(localhost:13316)/webook"
`
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		"http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置中心发生变更")
	})
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("watch", viper.GetString("test.key"))
			//time.Sleep(time.Second)
		}
	}()
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}
