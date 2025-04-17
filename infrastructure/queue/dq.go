package queue

import (
	"context"
	"encoding/json"
	"mq/infrastructure/svc"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/zeromicro/go-queue/dq"
	"github.com/zeromicro/go-zero/core/logx"
)

// DQConfig DQ队列配置
type DQConfig struct {
	// 重试配置
	MaxRetries     int
	RetryInterval  time.Duration
	ReconnectDelay time.Duration
	MaxReconnects  int
	// 消息延迟时间
	DelayTime time.Duration
}

// DefaultDQConfig 返回默认 DQ 配置
func DefaultDQConfig(svcCtx *svc.ServiceContext) *DQConfig {
	return &DQConfig{
		MaxRetries:     3,
		RetryInterval:  time.Second * 5,
		ReconnectDelay: time.Second * 5,
		MaxReconnects:  5,
		DelayTime:      time.Second * 5,
	}
}

// DQQueue 实现
type DQQueue struct {
	config   *DQConfig
	producer dq.Producer
	consumer dq.Consumer
	mu       sync.RWMutex
	isClosed bool
}

// NewDQQueue 创建新的 DQQueue 实例
func NewDQQueue(svcCtx *svc.ServiceContext) *DQQueue {
	config := DefaultDQConfig(svcCtx)

	// 创建生产者
	producer := dq.NewProducer(svcCtx.GetConfig().DqConf.Beanstalks)

	// 创建消费者
	consumer := dq.NewConsumer(svcCtx.GetConfig().DqConf)

	return &DQQueue{
		config:   config,
		producer: producer,
		consumer: consumer,
	}
}

// Publish 发布消息
func (q *DQQueue) Publish(ctx context.Context, message *Message) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.isClosed {
		return ErrQueueClosed
	}

	msgBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for i := 0; i < q.config.MaxRetries; i++ {
		_, err := q.producer.Delay(msgBody, q.config.DelayTime)
		if err == nil {
			return nil
		}
		time.Sleep(q.config.RetryInterval)
	}

	return ErrPublishFailed
}

// Consume 消费消息
func (q *DQQueue) Consume(ctx context.Context, queueType []MessageType, handler func(message *Message) error) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.isClosed {
		return ErrQueueClosed
	}

	// 使用单个 goroutine 来处理消费
	go func() {
		for {
			if q.isClosed {
				logx.Info("Queue is closed, stopping consumer")
				return
			}

			logx.Info("Starting consumer...")

			// 直接调用消费方法
			q.consumer.Consume(func(body []byte) {
				// 解析消息类型和ID
				var msg Message

				if err := json.Unmarshal(body, &msg); err != nil {
					logx.Errorf("Failed to unmarshal message: %v", err)
					return
				}
				if !lo.Contains(queueType, msg.Type) {
					logx.Infof("Message type %v not in queue types %v", msg.Type, queueType)
					return
				}
				if err := handler(&msg); err != nil {
					logx.Errorf("Failed to handle message: %v", err)
				} else {
					logx.Info("Message handled successfully")
				}
			})

			// 如果发生错误，等待一段时间后重试
			logx.Info("Consumer finished, waiting for retry...")
			time.Sleep(q.config.RetryInterval)
		}
	}()

	return nil
}

// Close 关闭连接
func (q *DQQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isClosed {
		return nil
	}

	q.isClosed = true

	// 关闭生产者
	if q.producer != nil {
		q.producer.Close()
	}

	return nil
}

// 错误定义
var (
	ErrQueueClosed   = NewError("queue is closed")
	ErrPublishFailed = NewError("failed to publish message")
)

// Error 自定义错误类型
type Error struct {
	message string
}

// NewError 创建新的错误
func NewError(message string) error {
	return &Error{message: message}
}

func (e *Error) Error() string {
	return e.message
}
