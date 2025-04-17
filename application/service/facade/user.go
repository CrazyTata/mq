package facade

import (
	"context"
	"mq/application/service"
	"mq/domain/statistics"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFacade struct {
	patientService     *service.PatientService
	appointmentService *service.AppointmentService
	healthService      *service.HealthService
	operationService   *service.OperationService
	statisticsService  *service.StatisticsService
	userService        *service.UserService
	cacheTTL           time.Duration
}

func NewUserFacade(patientService *service.PatientService, appointmentService *service.AppointmentService, healthService *service.HealthService, operationService *service.OperationService, statisticsService *service.StatisticsService, userService *service.UserService) *UserFacade {
	return &UserFacade{
		patientService:     patientService,
		appointmentService: appointmentService,
		healthService:      healthService,
		operationService:   operationService,
		statisticsService:  statisticsService,
		userService:        userService,
		cacheTTL:           5 * time.Minute, // 缓存5分钟
	}
}

// GetStatistics 获取统计信息
func (s *UserFacade) UpdateStatistics(ctx context.Context) {
	// 获取所有用户
	users, err := s.userService.GetAllUsers(ctx)
	if err != nil {
		return
	}

	for _, user := range users {
		go s.UpdateUserStatistics(user)
	}
	return
}
func (s *UserFacade) UpdateUserStatistics(userId string) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	var err error
	// 并发获取各项统计数据
	var wg sync.WaitGroup
	var totalPatients, activePatients, todayAppointments, upcomingAppointments, healthRecords, recentOperations int64

	wg.Add(6)
	go func() {
		defer wg.Done()
		totalPatients, err = s.patientService.CountByUserID(ctx, userId)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		activePatients, err = s.patientService.CountActiveByUserID(ctx, userId)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		todayAppointments, err = s.appointmentService.CountTodayByUserID(ctx, userId)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		upcomingAppointments, err = s.appointmentService.CountUpcomingByUserID(ctx, userId)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		healthRecords, err = s.healthService.CountByUserID(ctx, userId)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		recentOperations, err = s.operationService.CountRecentByUserID(ctx, userId, 30)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}()

	wg.Wait()

	if err != nil {
		logger.Errorf("UpdateUserStatistics error: %v", err)
	}
	statisticsDo, err := s.statisticsService.GetByDateAndUserId(ctx, userId, time.Now().Format("2006-01-02"))
	if err != nil {
		logger.Errorf("UpdateUserStatistics error: %v", err)
		return
	}
	if statisticsDo == nil {
		statisticsDo = statistics.Create(userId, totalPatients, activePatients, todayAppointments, upcomingAppointments, healthRecords, recentOperations, time.Now().Format("2006-01-02"))
		_, err = s.statisticsService.Create(ctx, statisticsDo)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
		return
	}
	statisticsDo.Update(totalPatients, activePatients, todayAppointments, upcomingAppointments, healthRecords, recentOperations)
	if err = s.statisticsService.Update(ctx, statisticsDo); err != nil {
		logger.Errorf("UpdateUserStatistics error: %v", err)
	}
	return
}

func (s *UserFacade) UpdatePatientStatistics(userId string) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	totalPatients, err := s.patientService.CountByUserID(ctx, userId)
	if err != nil {
		logger.Errorf("UpdatePatientStatistics  patientService.CountByUserID error: %v", err)
	}
	activePatients, err := s.patientService.CountActiveByUserID(ctx, userId)
	if err != nil {
		logger.Errorf("UpdatePatientStatistics  patientService.CountActiveByUserID error: %v", err)
	}
	statisticsDo, err := s.statisticsService.GetByDateAndUserId(ctx, userId, time.Now().Format("2006-01-02"))
	if err != nil {
		logger.Errorf("UpdatePatientStatistics  statisticsService.GetByDateAndUserId error: %v", err)
		return
	}
	if statisticsDo == nil {
		statisticsDo = statistics.Create(userId, totalPatients, activePatients, 0, 0, 0, 0, time.Now().Format("2006-01-02"))
		_, err = s.statisticsService.Create(ctx, statisticsDo)
		if err != nil {
			logger.Errorf("UpdatePatientStatistics  statisticsService.Create error: %v", err)
		}
	} else {
		statisticsDo.Update(totalPatients, activePatients, 0, 0, 0, 0)
		if err = s.statisticsService.Update(ctx, statisticsDo); err != nil {
			logger.Errorf("UpdatePatientStatistics  statisticsService.Update error: %v", err)
		}
	}
}

func (s *UserFacade) UpdateAppointmentStatistics(userId string) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	todayAppointments, err := s.appointmentService.CountTodayByUserID(ctx, userId)
	if err != nil {
		logger.Errorf("UpdateAppointmentStatistics  appointmentService.CountTodayByUserID error: %v", err)
	}
	upcomingAppointments, err := s.appointmentService.CountUpcomingByUserID(ctx, userId)
	if err != nil {
		logger.Errorf("UpdateAppointmentStatistics  appointmentService.CountUpcomingByUserID error: %v", err)
	}
	statisticsDo, err := s.statisticsService.GetByDateAndUserId(ctx, userId, time.Now().Format("2006-01-02"))
	if err != nil {
		logger.Errorf("UpdateAppointmentStatistics  statisticsService.GetByDateAndUserId error: %v", err)
		return
	}
	if statisticsDo == nil {
		statisticsDo = statistics.Create(userId, 0, 0, todayAppointments, upcomingAppointments, 0, 0, time.Now().Format("2006-01-02"))
		_, err = s.statisticsService.Create(ctx, statisticsDo)
		if err != nil {
			logger.Errorf("UpdateAppointmentStatistics  statisticsService.Create error: %v", err)
		}
		return
	}
	statisticsDo.Update(0, 0, todayAppointments, upcomingAppointments, 0, 0)
	if err = s.statisticsService.Update(ctx, statisticsDo); err != nil {
		logger.Errorf("UpdateUserStatistics error: %v", err)
	}
	return

}

func (s *UserFacade) UpdateHealthStatistics(userId string) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	healthRecords, err := s.healthService.CountByUserID(ctx, userId)
	if err != nil {
		logger.Errorf("UpdateHealthStatistics  healthService.CountByUserID error: %v", err)
	}
	statisticsDo, err := s.statisticsService.GetByDateAndUserId(ctx, userId, time.Now().Format("2006-01-02"))
	if err != nil {
		logger.Errorf("UpdateHealthStatistics  statisticsService.GetByDateAndUserId error: %v", err)
		return
	}
	if statisticsDo == nil {
		statisticsDo = statistics.Create(userId, 0, 0, 0, 0, healthRecords, 0, time.Now().Format("2006-01-02"))
		_, err = s.statisticsService.Create(ctx, statisticsDo)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	} else {
		statisticsDo.Update(0, 0, 0, 0, healthRecords, 0)
		if err = s.statisticsService.Update(ctx, statisticsDo); err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
	}
}

func (s *UserFacade) UpdateOperationStatistics(userId string) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	recentOperations, err := s.operationService.CountRecentByUserID(ctx, userId, 30)
	if err != nil {
		logger.Errorf("UpdateOperationStatistics  operationService.CountRecentByUserID error: %v", err)
	}

	statisticsDo, err := s.statisticsService.GetByDateAndUserId(ctx, userId, time.Now().Format("2006-01-02"))
	if err != nil {
		logger.Errorf("UpdateOperationStatistics  statisticsService.GetByDateAndUserId error: %v", err)
		return
	}
	if statisticsDo == nil {
		statisticsDo = statistics.Create(userId, 0, 0, 0, 0, 0, recentOperations, time.Now().Format("2006-01-02"))
		_, err = s.statisticsService.Create(ctx, statisticsDo)
		if err != nil {
			logger.Errorf("UpdateUserStatistics error: %v", err)
		}
		return
	}
	statisticsDo.Update(0, 0, 0, 0, 0, recentOperations)
	if err = s.statisticsService.Update(ctx, statisticsDo); err != nil {
		logger.Errorf("UpdateOperationStatistics  statisticsService.Update error: %v", err)
	}
	return
}
