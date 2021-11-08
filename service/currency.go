package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/mhaqiw/game-currency/domain"
	"time"
)

type currencyService struct {
	currencyRepo   domain.CurrencyRepository
	contextTimeout time.Duration
}

func NewCurrencyService(c domain.CurrencyRepository, timeout time.Duration) domain.CurrencyService {
	return &currencyService{
		currencyRepo:   c,
		contextTimeout: timeout,
	}
}

func (s currencyService) GetAll(ctx context.Context) (response domain.CurrenciesResponsePayload, err error) {
	list, err := s.currencyRepo.GetAll(ctx)
	if err != nil {
		return
	}

	response.List = []domain.CurrencyResponse{}
	err = copier.Copy(&response.List, &list)
	return

}

func (s currencyService) Create(ctx context.Context, request domain.CurrencyRequestPayload) ( currency domain.Currency,err error) {
	currency = domain.Currency{
		Name:      request.Name,
	}

	isAlreadyExists, err := s.currencyRepo.CheckIsExistByName(ctx, currency.Name)
	if err != nil {
		return
	}
	if isAlreadyExists {
		return currency, domain.ErrCurrencyAlreadyExists
	}
	err = s.currencyRepo.Create(ctx, &currency)
	return
}
