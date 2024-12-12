package saramax

import (
	"encoding/json"
	"geektime-basic-go/webook/pkg/logger"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	l  logger.LoggerV1
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](l logger.LoggerV1, fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for message := range messages {
		var t T
		err := json.Unmarshal(message.Value, &t)
		if err != nil {
			h.l.Error("反序列化消息失败", logger.Error(err),
				logger.String("topic", message.Topic),
				logger.Int32("partition", message.Partition),
				logger.Int64("offset", message.Offset))
			continue
		}
		err = h.fn(message, t)

		for i := 0; i < 3; i++ {
			err = h.fn(message, t)
			if err == nil {
				break
			}
			h.l.Error("处理消息失败", logger.Error(err),
				logger.String("topic", message.Topic),
				logger.Int32("partition", message.Partition),
				logger.Int64("offset", message.Offset))
		}
		// 在这里执行重试
		if err != nil {
			h.l.Error("处理消息失败-重试次数上限", logger.Error(err),
				logger.String("topic", message.Topic),
				logger.Int32("partition", message.Partition),
				logger.Int64("offset", message.Offset))
		} else {
			session.MarkMessage(message, "")
		}
	}
	return nil
}
