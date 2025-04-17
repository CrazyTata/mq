package queue

import (
	"context"
	"sync"
)

// ConsumerManager 消费者管理器
type ConsumerManager struct {
	consumers []Consumer
	queue     QueueManager
}

// Consumer 消费者接口
type Consumer interface {
	Start(ctx context.Context) error
}

// NewConsumerManager 创建消费者管理器
func NewConsumerManager(queue QueueManager) *ConsumerManager {
	return &ConsumerManager{
		queue: queue,
	}
}

// Register 注册消费者
func (m *ConsumerManager) Register(consumer Consumer) {
	m.consumers = append(m.consumers, consumer)
}

// StartAll 启动所有消费者
func (m *ConsumerManager) StartAll(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(m.consumers))

	for _, consumer := range m.consumers {
		wg.Add(1)
		go func(c Consumer) {
			defer wg.Done()
			if err := c.Start(ctx); err != nil {
				errChan <- err
			}
		}(consumer)
	}

	// 等待所有消费者启动完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
