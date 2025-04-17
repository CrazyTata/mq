package service

import (
	"context"
	"errors"
	"fmt"
	"mq/application/assembler"
	"mq/application/dto"
	"mq/common/util"
	"mq/domain/patient"
)

// PatientService 患者应用服务
type PatientService struct {
	repo patient.PatientRepository
}

// NewPatientAppService 创建患者应用服务
func NewPatientService(repo patient.PatientRepository) *PatientService {
	return &PatientService{
		repo: repo,
	}
}

// Save 创建患者
func (s *PatientService) Save(ctx context.Context, req *dto.Patient) (*dto.CreatePatientResponse, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var patientID int64
	if req.Id == 0 {
		patient := patient.Create(req.Name, req.Age, req.Gender, req.Phone, userID, req.Status, req.Avatar, req.History, req.Allergies, req.Note, req.Attachments, req.Details)
		patientID, err = s.repo.Create(ctx, patient)
		if err != nil {
			return nil, err
		}
	} else {
		patient, err := s.repo.GetByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		patient.UpdateProfile(req.Name, req.Age, req.Gender, req.Phone, req.Avatar, req.History, req.Allergies, req.Note, req.Attachments, req.Details)
		err = s.repo.UpdatePatient(ctx, patient)
		if err != nil {
			return nil, err
		}
	}

	return &dto.CreatePatientResponse{
		ID: patientID,
	}, nil
}

// GetPatientByID 根据ID获取患者
func (s *PatientService) GetPatientByID(ctx context.Context, id int64) (*dto.Patient, error) {
	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOPatient(patient), nil
}

// GetPatientByFriendlyID 根据友好ID获取患者
func (s *PatientService) GetPatientByFriendlyID(ctx context.Context, friendlyID string) (*dto.Patient, error) {
	patient, err := s.repo.GetByFriendlyID(ctx, friendlyID)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOPatient(patient), nil
}

// GetPatientList 获取患者列表
func (s *PatientService) GetPatientList(ctx context.Context, req *dto.GetPatientListRequest) (*dto.PatientListResponse, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	patients, total, err := s.repo.GetList(ctx, userID, req.Search, req.Status, req.Order, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOPatients(patients, total), nil
}

// DeletePatient 删除患者
func (s *PatientService) DeletePatient(ctx context.Context, id int64) error {
	//先查询是否存在
	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if patient == nil {
		return errors.New("patient not found")
	}
	fmt.Println("patient", patient)
	patient.MarkAsDeleted()
	//删除患者
	return s.repo.UpdatePatient(ctx, patient)
}

// CountByUserID 根据用户ID获取患者数量
func (s *PatientService) CountByUserID(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountByUserID(ctx, userID)
}

// CountActiveByUserID 根据用户ID获取活跃患者数量
func (s *PatientService) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountActiveByUserID(ctx, userID)
}
