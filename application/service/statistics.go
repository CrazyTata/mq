package service

import (
	"context"
	"mq/application/assembler"
	"mq/application/dto"
	"mq/common/util"
	"mq/domain/statistics"
	"time"
)

type StatisticsService struct {
	statisticsRepo statistics.StatisticsRepository
}

func NewStatisticsService(statisticsRepo statistics.StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepo: statisticsRepo,
	}
}

func (s *StatisticsService) GetByDateAndUserId(ctx context.Context, userId, date string) (*statistics.Statistics, error) {
	return s.statisticsRepo.GetStatistics(ctx, userId, date)
}

func (s *StatisticsService) Create(ctx context.Context, statistics *statistics.Statistics) (int64, error) {
	return s.statisticsRepo.InsertStatistics(ctx, statistics)
}

func (s *StatisticsService) Update(ctx context.Context, statistics *statistics.Statistics) error {
	return s.statisticsRepo.UpdateStatistics(ctx, statistics)
}

func (s *StatisticsService) GetNowStatistics(ctx context.Context) (*dto.StatisticsResponse, error) {
	userId, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	statistics, err := s.statisticsRepo.GetStatistics(ctx, userId, time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	return assembler.DOTODTOStatistics(statistics), nil
}
