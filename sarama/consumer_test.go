package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

var addrs1 = []string{"localhost:9094"}

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	// 正常来说，一个消费者都是归属于一个消费者组的
	// 消费者组就是你的业务
	consumer, err := sarama.NewConsumerGroup(addrs1, "test_group", cfg)
	require.NoError(t, err)

	// 带超时的context
	start := time.Now()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Minute*5, func() {
		cancel()
	})
	err = consumer.Consume(ctx, []string{"test_topic"}, testConsumerGroupHandler{})
	// 消费结束，到这里
	t.Log(err, time.Since(start).String())
}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// topic => 偏移量
	partitions := session.Claims()["test_topic"]
	for _, part := range partitions {
		session.ResetOffset("test_topic", part, sarama.OffsetOldest, "")
	}
	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

func (t testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()

	var batchSize = 10

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var eg errgroup.Group
		var last *sarama.ConsumerMessage
		for i := 0; i < batchSize; i++ {
			done := false
			select {
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					// 消费者关闭
					return nil
				}
				last = msg
				eg.Go(func() error {
					// 在这里消费
					time.Sleep(time.Second)
					// 在这重试
					log.Println(string(msg.Value))
					return nil
				})
			case <-ctx.Done():
				// 这边代表超时
				break
			}
			if done {
				break
			}

		}
		cancel()
		err := eg.Wait()
		if err != nil {
			// 这边能怎么办
			// 记录日志
			continue
		}

		if last != nil {
			session.MarkMessage(last, "")
		}

	}

	return nil
}

func (t testConsumerGroupHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		/*var bizMsg MyBizMsg
		err := json.Unmarshal(msg.Value, &bizMsg)
		if err != nil {
			// 消费出错
			// 重试
			continue
		}*/
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

type MyBizMsg struct {
	Name string
}

// ChannelV1 返回一个只读的channel
func ChannelV1() <-chan struct{} {
	panic("implement me")
}

// ChannelV2 返回可读可写的channel
func ChannelV2() chan struct{} {
	panic("implement me")
}

// ChannelV3 返回只写的channel
func ChannelV3() chan<- struct{} {
	panic("implement me")
}
