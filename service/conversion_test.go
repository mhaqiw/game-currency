package service

import (
	"context"
	"github.com/mhaqiw/game-currency/domain"
	"github.com/mhaqiw/game-currency/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

func Test_conversionService_Create(t *testing.T) {
	timeNow := time.Now()
	payload := domain.ConversionInsertRequestPayload{
		FromID: 1,
		ToID:   2,
		Rate:   10,
	}
	tests := []struct {
		name           string
		request        domain.ConversionInsertRequestPayload
		conversionRepo mocks.ConversionRepository
		currencyRepo   mocks.CurrencyRepository
		wantResponse   domain.ConversionInsertResponsePayload
		wantErr        bool
	}{
		{
			name:    "success",
			request: payload,
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				repo.On("Post", mock.Anything, mock.Anything).Return(int64(1), int64(2), timeNow, nil).Once()
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(true, nil)
				return *repo
			}(),
			wantResponse: domain.ConversionInsertResponsePayload{List: []domain.Conversion{
				{ID: 1, FromID: 1, ToID: 2, Rate: 10, CreatedAt: timeNow}, {ID: 2, FromID: 2, ToID: 1, Rate: 0.1, CreatedAt: timeNow},
			}},
			wantErr: false,
		},
		{
			name:    "error insert",
			request: payload,
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				repo.On("Post", mock.Anything, mock.Anything).Return(int64(0), int64(0), timeNow, domain.ErrSqlError).Once()
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(true, nil)
				return *repo
			}(),
			wantResponse: domain.ConversionInsertResponsePayload{},
			wantErr:      true,
		},
		{
			name:    "error from_id not exist",
			request: payload,
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(false, nil)
				return *repo
			}(),
			wantResponse: domain.ConversionInsertResponsePayload{},
			wantErr:      true,
		},
		{
			name:    "error check from_id",
			request: payload,
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(false, domain.ErrSqlError)
				return *repo
			}(),
			wantResponse: domain.ConversionInsertResponsePayload{},
			wantErr:      true,
		},
		{
			name:    "error to_id not exist",
			request: payload,
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(true, nil).Once()
				repo.On("CheckIsExistByID", mock.Anything, int64(2)).Return(false, nil).Once()
				return *repo
			}(),
			wantResponse: domain.ConversionInsertResponsePayload{},
			wantErr:      true,
		},
		{
			name:    "error check to_id ",
			request: payload,
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(true, nil).Once()
				repo.On("CheckIsExistByID", mock.Anything, int64(2)).Return(false, domain.ErrSqlError).Once()
				return *repo
			}(),
			wantResponse: domain.ConversionInsertResponsePayload{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewConversionService(&tt.conversionRepo, &tt.currencyRepo, 5)
			gotResponse, err := s.Create(context.TODO(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Create() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_conversionService_GetConvertedCurrency(t *testing.T) {
	payload := domain.ConvertedCurrencyRequest{
		FromID: 1,
		ToID:   2,
		Amount: 10,
	}
	tests := []struct {
		name           string
		conversionRepo mocks.ConversionRepository
		currencyRepo   mocks.CurrencyRepository
		request        domain.ConvertedCurrencyRequest
		wantResponse   domain.ConvertedCurrencyResponse
		wantErr        bool
	}{
		{
			name: "success",
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				repo.On("GetFromToID", mock.Anything, mock.Anything, mock.Anything).Return(domain.Conversion{Rate: 10}, nil)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(true, nil)
				return *repo
			}(),
			request: payload,
			wantResponse: domain.ConvertedCurrencyResponse{
				Result: 100,
			},
			wantErr: false,
		},
		{
			name: "error get conversion",
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				repo.On("GetFromToID", mock.Anything, mock.Anything, mock.Anything).Return(domain.Conversion{}, domain.ErrSqlError)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(true, nil)
				return *repo
			}(),
			request:      payload,
			wantResponse: domain.ConvertedCurrencyResponse{},
			wantErr:      true,
		},
		{
			name: "from_id not found",
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(false, nil)
				return *repo
			}(),
			request:      payload,
			wantResponse: domain.ConvertedCurrencyResponse{},
			wantErr:      true,
		},
		{
			name: "error check from_id",
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, mock.Anything).Return(false, domain.ErrSqlError)
				return *repo
			}(),
			request:      payload,
			wantResponse: domain.ConvertedCurrencyResponse{},
			wantErr:      true,
		},
		{
			name: "to_id not found",
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, int64(1)).Return(true, nil).Once()
				repo.On("CheckIsExistByID", mock.Anything, int64(2)).Return(false, nil).Once()
				return *repo
			}(),
			request:      payload,
			wantResponse: domain.ConvertedCurrencyResponse{},
			wantErr:      true,
		},
		{
			name: "error check to_id",
			conversionRepo: func() mocks.ConversionRepository {
				repo := new(mocks.ConversionRepository)
				return *repo
			}(),
			currencyRepo: func() mocks.CurrencyRepository {
				repo := new(mocks.CurrencyRepository)
				repo.On("CheckIsExistByID", mock.Anything, int64(1)).Return(true, nil).Once()
				repo.On("CheckIsExistByID", mock.Anything, int64(2)).Return(false, domain.ErrSqlError).Once()
				return *repo
			}(),
			request:      payload,
			wantResponse: domain.ConvertedCurrencyResponse{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewConversionService(&tt.conversionRepo, &tt.currencyRepo, 5)
			gotResponse, err := s.GetConvertedCurrency(context.TODO(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConvertedCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("GetConvertedCurrency() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
