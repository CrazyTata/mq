package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"mq/application/service/facade"
	"mq/domain/health"
	"mq/infrastructure/queue"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthConsumer struct {
	queue      queue.QueueManager
	userFacade *facade.UserFacade
}

func NewHealthConsumer(queue queue.QueueManager, userFacade *facade.UserFacade) *HealthConsumer {
	return &HealthConsumer{queue: queue, userFacade: userFacade}
}

func (c *HealthConsumer) Start(ctx context.Context) error {
	return c.queue.Consume(ctx, []queue.MessageType{queue.HealthRecordSaved}, c.handleMessage)
}

func (c *HealthConsumer) handleMessage(message *queue.Message) error {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	if message == nil {
		logger.Errorf("message is nil")
		return nil
	}

	var healthRecord health.HealthRecords
	if err := json.Unmarshal(message.Body, &healthRecord); err != nil {
		logger.Errorf("failed to unmarshal health record: %v", err)
		return fmt.Errorf("failed to unmarshal health record: %v", err)
	}
	c.userFacade.UpdateHealthStatistics(healthRecord.UserId)
	return nil
}
