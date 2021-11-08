package service

import (
	"context"
	"github.com/mhaqiw/game-currency/domain"
	"time"
)

type conversionService struct {
	conversionRepo domain.ConversionRepository
	currencyRepo   domain.CurrencyRepository
	contextTimeout time.Duration
}

func NewConversionService(c domain.ConversionRepository, cu domain.CurrencyRepository, timeout time.Duration) domain.ConversionService {
	return &conversionService{
		conversionRepo: c,
		currencyRepo:   cu,
		contextTimeout: timeout,
	}
}

func (s conversionService) Create(ctx context.Context, request domain.ConversionInsertRequestPayload) (response domain.ConversionInsertResponsePayload, err error) {
	err = s.checkCurrencyFromAndTo(ctx, request.FromID, request.ToID)
	if err != nil {
		return
	}

	conversion1 := domain.Conversion{
		FromID: request.FromID,
		ToID:   request.ToID,
		Rate:   request.Rate,
	}

	conversion2 := domain.Conversion{
		FromID: request.ToID,
		ToID:   request.FromID,
		Rate:   1 / request.Rate,
	}

	conversion1ID, conversion2ID, createdAt, err := s.conversionRepo.Post(ctx, conversion1)
	if err != nil {
		return
	}

	conversion1.ID = conversion1ID
	conversion1.CreatedAt = createdAt
	conversion2.ID = conversion2ID
	conversion2.CreatedAt = createdAt
	response.List = []domain.Conversion{
		conversion1, conversion2,
	}

	return
}

func (s conversionService) GetConvertedCurrency(ctx context.Context, request domain.ConvertedCurrencyRequest) (response domain.ConvertedCurrencyResponse, err error) {
	err = s.checkCurrencyFromAndTo(ctx, request.FromID, request.ToID)
	if err != nil {
		return
	}

	conversion, err := s.conversionRepo.GetFromToID(ctx, request.FromID, request.ToID)
	if err != nil {
		return
	}

	response.Result = conversion.Rate * request.Amount
	return
}

func (s conversionService) checkCurrencyFromAndTo(ctx context.Context, fromID int64, toID int64) (err error) {
	isFromIDExist, err := s.currencyRepo.CheckIsExistByID(ctx, fromID)
	if err != nil {
		return
	}
	if !isFromIDExist {
		return domain.ErrCurrencyFromIDNotExists
	}

	isToIDExist, err := s.currencyRepo.CheckIsExistByID(ctx, toID)
	if err != nil {
		return
	}
	if !isToIDExist {
		return domain.ErrCurrencyToIDNotExists
	}

	return nil
}
