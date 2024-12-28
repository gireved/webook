package ioc

import (
	"geektime-basic-go/webook/internal/repository/dao"
	"geektime-basic-go/webook/pkg/logger"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	"time"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:3316)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(goormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	// 监控查询的执行时间
	pcb := newCallbacks()
	//pcb.registerAll(db)
	db.Use(pcb)
	return db

}

type goormLoggerFunc func(msg string, fields ...logger.Field)

func (g goormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{Key: "args", Val: i})
}

type Callbacks struct {
	vector *prometheus2.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus-query"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	c.registerAll(db)
	return nil
}

func newCallbacks() *Callbacks {
	vector := prometheus2.NewSummaryVec(prometheus2.SummaryOpts{
		// 在这边考虑 设置各种namespace
		Namespace: "webook_gorm_whx",
		Subsystem: "webook",
		Name:      "gorm_query_time",
		Help:      "统计 GORM 的执行时间",
		ConstLabels: map[string]string{
			"db": "webook",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})

	pcb := &Callbacks{
		vector: vector,
	}
	prometheus2.MustRegister(vector)
	return pcb
}

func (c *Callbacks) registerAll(db *gorm.DB) {
	// 作用于 insert 语句
	err := db.Callback().Create().Before("*").Register("prometheus_create_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_after", c.after("create"))
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	err = db.Callback().Update().Before("*").Register("prometheus_update_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").Register("prometheus_update_after", c.after("update"))
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	err = db.Callback().Delete().Before("*").Register("prometheus_delete_after", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").Register("prometheus_delete_after", c.after("delete"))
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	err = db.Callback().Raw().Before("*").Register("prometheus_raw_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").Register("prometheus_raw_after", c.after("raw"))
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	err = db.Callback().Row().Before("*").Register("prometheus_row_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").Register("prometheus_row_after", c.after("row"))
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

}

func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}

func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			// 你啥都干不了
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		// 准备上报prometheus
		c.vector.WithLabelValues(typ, table).Observe(float64(time.Since(startTime).Milliseconds()))

	}
}
