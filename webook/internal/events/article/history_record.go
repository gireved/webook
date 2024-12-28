package article

import (
	"context"
	"geektime-basic-go/webook/internal/repository"
	"geektime-basic-go/webook/pkg/logger"
	"geektime-basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type HistoryRecordEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewHistoryRecordEventConsumer(client sarama.Client,
	l logger.LoggerV1, repo repository.InteractiveRepository) *HistoryRecordEventConsumer {
	return &HistoryRecordEventConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (i *HistoryRecordEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (i *HistoryRecordEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.AddRecord(ctx, t.Aid, t.Uid)
}
