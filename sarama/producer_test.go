package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	//cfg.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		// 消息数据本体
		// 转JSON
		Value: sarama.StringEncoder(`{"aid":1,"uid":123}`),
		// 会在生产者和消费者之间传递
		/*	Headers: []sarama.RecordHeader{
				{
					Key:   []byte("trace_id"),
					Value: []byte("123456"),
				},
			},
			// 只作用于发送过程
			Metadata: "这是metadata",*/
	})
	assert.NoError(t, err)
}

type JSONEncoder struct {
	Data any
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)
	msgCh := producer.Input()
	go func() {
		for i := 0; i < 10; i++ {
			msgCh <- &sarama.ProducerMessage{
				Topic: "test_topic",
				Key:   sarama.StringEncoder("oid-123"),
				// 消息数据本体
				// 转JSON
				Value: sarama.StringEncoder("hello 这是一条消息 B"),
				// 会在生产者和消费者之间传递
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("trace_id"),
						Value: []byte("123456"),
					},
				},
				// 只作用于发送过程
				Metadata: "这是metadata",
			}
		}
	}()

	errCh := producer.Errors()
	succCh := producer.Successes()

	select {
	case err := <-errCh:
		t.Log("发送出了问题", err.Err)
	case <-succCh:
		t.Log("发送成功")
	}
}
