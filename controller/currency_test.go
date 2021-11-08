package controller

import (
	"github.com/labstack/echo"
	"github.com/mhaqiw/game-currency/domain"
	"github.com/mhaqiw/game-currency/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCurrencyHandler_GetAll(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name                string
		mockCurrencyService mocks.CurrencyService
		statusCodeExpected  int
		wantErr             bool
		responseBody        string
	}{
		{
			name: "success",
			mockCurrencyService: func() mocks.CurrencyService {
				service := new(mocks.CurrencyService)
				service.On("GetAll", mock.Anything).Return(domain.CurrenciesResponsePayload{
					List: []domain.CurrencyResponse{
						{ID: 1, Name: "Knut"},
					},
				}, nil)
				return *service
			}(),
			statusCodeExpected: http.StatusOK,
			wantErr:            false,
			responseBody:       "{\"list\":[{\"id\":1,\"name\":\"Knut\"}]}\n",
		},
		{
			name: "error",
			mockCurrencyService: func() mocks.CurrencyService {
				service := new(mocks.CurrencyService)
				service.On("GetAll", mock.Anything).Return(domain.CurrenciesResponsePayload{}, domain.ErrSqlError)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			wantErr:            true,
			responseBody:       "{\"message\":\"sql error\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewCurrencyHandler(e, &tt.mockCurrencyService)
			req := httptest.NewRequest(http.MethodGet, "/currency", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, tt.statusCodeExpected, rec.Code)
			require.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}

func TestCurrencyHandler_Post(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name                string
		mockCurrencyService mocks.CurrencyService
		statusCodeExpected  int
		wantErr             bool
		responseBody        string
		reqBody             io.Reader
	}{
		{
			name: "success",
			mockCurrencyService: func() mocks.CurrencyService {
				service := new(mocks.CurrencyService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.Currency{}, nil)
				return *service
			}(),
			statusCodeExpected: http.StatusCreated,
			wantErr:            false,
			reqBody: strings.NewReader(`{"name": "cukong"}`),
			responseBody: "{\"id\":0,\"name\":\"\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
		},
		{
			name: "error sql",
			mockCurrencyService: func() mocks.CurrencyService {
				service := new(mocks.CurrencyService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.Currency{}, domain.ErrSqlError)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			wantErr:            true,
			reqBody: strings.NewReader(`{"name": "cukong"}`),
			responseBody: "{\"message\":\"sql error\"}\n",
		},
		{
			name: "error invalid body request",
			mockCurrencyService: func() mocks.CurrencyService {
				service := new(mocks.CurrencyService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.Currency{}, domain.ErrSqlError)
				return *service
			}(),
			statusCodeExpected: http.StatusUnprocessableEntity,
			wantErr:            true,
			responseBody:       "{\"message\":\"code=400, message=Request body can't be empty\"}\n",
		},
		{
			name: "error empty name",
			mockCurrencyService: func() mocks.CurrencyService {
				service := new(mocks.CurrencyService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.Currency{}, domain.ErrSqlError)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			wantErr:            true,
			reqBody: strings.NewReader(`{"name": ""}`),
			responseBody: "{\"message\":\"Key: 'CurrencyRequestPayload.Name' Error:Field validation for 'Name' failed on the 'required' tag\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewCurrencyHandler(e, &tt.mockCurrencyService)
			req := httptest.NewRequest(http.MethodPost, "/currency", tt.reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, tt.statusCodeExpected, rec.Code)
			require.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}
