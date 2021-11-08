package service

import (
	"context"
	"github.com/mhaqiw/game-currency/domain"
	"github.com/mhaqiw/game-currency/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func Test_currencyService_Create(t *testing.T) {
	tests := []struct {
		name             string
		mockCurrencyRepo mocks.CurrencyRepository
		request          domain.CurrencyRequestPayload
		wantCurrency     domain.Currency
		wantErr          bool
	}{
		{
			name: "sucess",
			mockCurrencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByName", mock.Anything, "Knut").Return(false, nil).Once()
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
				return *repo
			}(),
			request: domain.CurrencyRequestPayload{
				Name: "Knut",
			},
			wantCurrency: domain.Currency{
				Name: "Knut",
			},
			wantErr: false,
		},
		{
			name: "error check exists",
			mockCurrencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByName", mock.Anything, "Knut").Return(false, domain.ErrSqlError).Once()
				return *repo
			}(),
			request: domain.CurrencyRequestPayload{
				Name: "Knut",
			},
			wantCurrency: domain.Currency{
				Name: "Knut",
			},
			wantErr: true,
		},
		{
			name: "error insert ",
			mockCurrencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByName", mock.Anything, "Knut").Return(false, nil).Once()
				repo.On("Create", mock.Anything, mock.Anything).Return(domain.ErrSqlError)
				return *repo
			}(),
			request: domain.CurrencyRequestPayload{
				Name: "Knut",
			},
			wantCurrency: domain.Currency{
				Name: "Knut",
			},
			wantErr: true,
		},
		{
			name: "failed already exists",
			mockCurrencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByName", mock.Anything, "Knut").Return(true, nil).Once()
				return *repo
			}(),
			request: domain.CurrencyRequestPayload{
				Name: "Knut",
			},
			wantCurrency: domain.Currency{
				Name: "Knut",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCurrencyService(&tt.mockCurrencyRepo, 5)
			gotCurrency, err := s.Create(context.TODO(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCurrency, tt.wantCurrency) {
				t.Errorf("Create() gotCurrency = %v, want %v", gotCurrency, tt.wantCurrency)
			}
		})
	}
}

func Test_currencyService_GetAll(t *testing.T) {
	tests := []struct {
		name             string
		mockCurrencyRepo mocks.CurrencyRepository
		wantResponse     domain.CurrenciesResponsePayload
		wantErr          bool
	}{
		{
			name: "success",
			mockCurrencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("GetAll", mock.Anything).Return([]domain.Currency{
					{
						Name: "Knut",
					},
				}, nil).Once()
				return *repo
			}(),
			wantResponse: domain.CurrenciesResponsePayload{List: []domain.CurrencyResponse{{
				Name: "Knut",
			},
			}},
			wantErr: false,
		},
		{
			name: "error",
			mockCurrencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("GetAll", mock.Anything).Return([]domain.Currency{}, domain.ErrSqlError).Once()
				return *repo
			}(),
			wantResponse: domain.CurrenciesResponsePayload{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCurrencyService(&tt.mockCurrencyRepo, 5)
			gotResponse, err := s.GetAll(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("GetAll() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
