package operation_records

import (
	"mq/domain/operation"
)

func POTODOGetOperationRecords(res *OperationRecords) *operation.OperationRecords {
	return &operation.OperationRecords{
		Id:        res.Id,
		Action:    res.Action,
		Target:    res.Target,
		Details:   res.Details,
		Username:  res.Username,
		UserId:    res.UserId,
		IsDeleted: res.IsDeleted,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}
}

func DOTOPOOperationRecords(res *operation.OperationRecords) *OperationRecords {
	return &OperationRecords{
		Id:        res.Id,
		Action:    res.Action,
		Target:    res.Target,
		Details:   res.Details,
		Username:  res.Username,
		UserId:    res.UserId,
		IsDeleted: res.IsDeleted,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}
}

func POTODOGetOperationRecordsList(res []OperationRecords) []*operation.OperationRecords {
	resp := make([]*operation.OperationRecords, 0, len(res))
	for _, v := range res {
		resp = append(resp, POTODOGetOperationRecords(&v))
	}
	return resp
}
