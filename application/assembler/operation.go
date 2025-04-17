package assembler

import (
	"mq/application/dto"
	"mq/domain/operation"
	"time"
)

func DOTODTOOperationList(operations []*operation.OperationRecords, total int64) *dto.GetOperationResponse {
	data := make([]*dto.OperationResponse, 0)
	for _, operation := range operations {
		data = append(data, DOTODTOOperation(operation))
	}
	return &dto.GetOperationResponse{
		Count: total,
		Data:  data,
	}
}

func DOTODTOOperation(operation *operation.OperationRecords) *dto.OperationResponse {
	if operation == nil {
		return nil
	}
	return &dto.OperationResponse{
		ID:        operation.Id,
		Action:    operation.Action,
		Target:    operation.Target,
		Details:   operation.Details.String,
		Username:  operation.Username,
		UserId:    operation.UserId,
		CreatedAt: operation.CreatedAt.Format(time.DateTime),
	}
}
