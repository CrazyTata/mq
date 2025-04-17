package service

import (
	"context"
	"errors"
	"mq/application/assembler"
	"mq/application/dto"
	"mq/common/util"
	"mq/domain/appointment"
)

// AppointmentService 预约应用服务
type AppointmentService struct {
	repo appointment.AppointmentRepository
}

// NewAppointmentAppService 创建预约应用服务
func NewAppointmentService(repo appointment.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		repo: repo,
	}
}

// Save 创建预约
func (s *AppointmentService) Save(ctx context.Context, req *dto.AppointmentRequest) (*dto.CreateAppointmentResponse, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var appointmentID int64
	if req.Id == 0 {
		appointment := appointment.Create(req.PatientID, req.PatientName, req.Date, req.Time, req.Duration, req.Type, req.Status, req.Notes, userID)
		appointmentID, err = s.repo.Create(ctx, appointment)
		if err != nil {
			return nil, err
		}
	} else {
		appointment, err := s.repo.GetByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		if appointment == nil || appointment.Id == 0 {
			return nil, errors.New("appointment not found")
		}
		appointment.Update(req.PatientName, req.Date, req.Time, req.Duration, req.Type, req.Status, req.Notes)
		err = s.repo.UpdateAppointment(ctx, appointment)
		if err != nil {
			return nil, err
		}
	}

	return &dto.CreateAppointmentResponse{
		ID: appointmentID,
	}, nil
}

// GetByID 根据ID获取预约
func (s *AppointmentService) GetByID(ctx context.Context, id int64) (*dto.AppointmentResponse, error) {
	appointment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOAppointment(appointment), nil
}

// Delete 删除预约
func (s *AppointmentService) Delete(ctx context.Context, id int64) error {
	appointment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if appointment == nil || appointment.Id == 0 {
		return errors.New("appointment not found")
	}
	appointment.MarkAsDeleted()
	return s.repo.UpdateAppointment(ctx, appointment)
}

// GetList 获取预约列表
func (s *AppointmentService) GetList(ctx context.Context, req *dto.GetAppointmentListRequest) (*dto.AppointmentListResponse, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	appointments, total, err := s.repo.GetList(ctx, userID, req.Search, req.SearchType, req.Order, req.PatientId, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOAppointmentList(appointments, total), nil
}

// CountTodayByUserID 根据用户ID获取今日预约数量
func (s *AppointmentService) CountTodayByUserID(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountTodayByUserID(ctx, userID)
}

// CountUpcomingByUserID 根据用户ID获取即将进行的预约数量
func (s *AppointmentService) CountUpcomingByUserID(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountUpcomingByUserID(ctx, userID)
}
