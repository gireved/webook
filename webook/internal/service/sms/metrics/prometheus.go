package metrics

import (
	"context"
	"geektime-basic-go/webook/internal/service/sms"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusDecorator struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewPrometheusDecorator(svc sms.Service) *PrometheusDecorator {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "webook_whx",
		Subsystem: "webook",
		Name:      "sms_resp_time",
		Help:      "统计 SMS 服务的性能",
	}, []string{"tplId"})

	prometheus.MustRegister(vector)

	return &PrometheusDecorator{
		svc:    svc,
		vector: vector,
	}
}

func (p *PrometheusDecorator) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		p.vector.WithLabelValues(tplId).Observe(float64(duration))
	}()
	return p.svc.Send(ctx, tplId, args, numbers...)

}
