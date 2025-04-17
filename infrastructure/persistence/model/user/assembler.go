package user

import (
	"mq/domain/user"
)

func POTODOGetUser(res *User) *user.User {
	return &user.User{
		ID:             res.Id,
		UserID:         res.UserId,
		Phone:          res.Phone,
		AppleID:        res.AppleId,
		Email:          res.Email,
		IsPrivateEmail: res.IsPrivateEmail,
		Name:           res.Name,
		Avatar:         res.Avatar,
		Password:       res.Password,
		RegisterTime:   res.RegisterTime,
		Status:         res.Status,
		Source:         res.Source,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}
}

func DOTOPOUser(res *user.User) *User {
	return &User{
		UserId:         res.UserID,
		Phone:          res.Phone,
		Password:       res.Password,
		Name:           res.Name,
		Avatar:         res.Avatar,
		Email:          res.Email,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
		Id:             res.ID,
		AppleId:        res.AppleID,
		IsPrivateEmail: res.IsPrivateEmail,
		RegisterTime:   res.RegisterTime,
		Status:         res.Status,
		Source:         res.Source,
	}
}

func POTODOGetUsers(res []User) []*user.User {
	resp := make([]*user.User, 0, len(res))
	for _, v := range res {
		resp = append(resp, POTODOGetUser(&v))
	}
	return resp
}
