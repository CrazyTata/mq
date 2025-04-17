package assembler

import (
	"mq/application/dto"
	domain "mq/domain/user"
	userModel "mq/infrastructure/persistence/model/user"
	"time"
)

func DOTOPOUser(do *domain.User) *userModel.User {
	return &userModel.User{
		Id:             do.ID,
		AppId:          do.AppId,
		UserId:         do.UserID,
		Phone:          do.Phone,
		AppleId:        do.AppleID,
		Email:          do.Email,
		IsPrivateEmail: do.IsPrivateEmail,
		Name:           do.Name,
		Avatar:         do.Avatar,
		Password:       do.Password,
		RegisterTime:   do.RegisterTime,
		Status:         do.Status,
		Source:         do.Source,
		CreatedAt:      do.CreatedAt,
		UpdatedAt:      do.UpdatedAt,
	}
}

func POTODOUser(po *userModel.User) *domain.User {
	return &domain.User{
		ID:             po.Id,
		AppId:          po.AppId,
		UserID:         po.UserId,
		Phone:          po.Phone,
		AppleID:        po.AppleId,
		Email:          po.Email,
		IsPrivateEmail: po.IsPrivateEmail,
		Name:           po.Name,
		Avatar:         po.Avatar,
		Password:       po.Password,
		RegisterTime:   po.RegisterTime,
		Status:         po.Status,
		Source:         po.Source,
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
	}
}

func DOTODTOUser(do *domain.User) *dto.GetUserResponse {
	if do == nil {
		return nil
	}
	registerTime := ""
	if do.RegisterTime.Valid {
		registerTime = do.RegisterTime.Time.Format(time.DateTime)
	}

	return &dto.GetUserResponse{
		ID:             do.ID,
		AppId:          do.AppId,
		UserID:         do.UserID,
		Name:           do.Name,
		Avatar:         do.Avatar,
		Phone:          do.GetDecryptedPhone(),
		Email:          do.Email,
		IsPrivateEmail: do.IsPrivateEmail,
		RegisterTime:   registerTime,
		Status:         do.Status,
		Source:         do.Source,
		CreatedAt:      do.CreatedAt.Format(time.DateTime),
	}
}

func DOTODTOBackendUser(users []*domain.User, total int64) *dto.BackendUserResponse {
	backendUsers := make([]*dto.BackendUser, 0, len(users))
	for _, u := range users {
		backendUser := &dto.BackendUser{
			UserID:    u.UserID,
			Phone:     u.GetDecryptedPhone(),
			AppleID:   u.AppleID,
			Email:     u.Email,
			Name:      u.Name,
			Avatar:    u.Avatar,
			Source:    u.Source,
			CreatedAt: u.CreatedAt.Format(time.DateTime),
		}

		// 处理可能为空的时间字段
		if u.RegisterTime.Valid {
			backendUser.RegisterTime = u.RegisterTime.Time.Format(time.DateTime)
		}

		backendUsers = append(backendUsers, backendUser)
	}

	return &dto.BackendUserResponse{
		Data:  backendUsers,
		Total: total,
	}
}
