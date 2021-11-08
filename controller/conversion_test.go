package controller

import (
	"fmt"
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

func TestConversionHandler_Post(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name                  string
		mockConversionService mocks.ConversionService
		statusCodeExpected    int
		wantErr               bool
		responseBody          string
		reqBody               io.Reader
	}{
		{
			name: "success",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.ConversionInsertResponsePayload{
					List: []domain.Conversion{{ID: 1, FromID: 1, ToID: 2, Rate: 10}, {ID: 2, FromID: 2, ToID: 1, Rate: 0.1}},
				}, nil)
				return *service
			}(),
			statusCodeExpected: http.StatusCreated,
			wantErr:            false,
			reqBody:            strings.NewReader(`{"from_id": 1,"to_id": 2,"rate": 29}`),
			responseBody:       "{\"list\":[{\"id\":1,\"from_id\":1,\"to_to\":2,\"rate\":10,\"created_at\":\"0001-01-01T00:00:00Z\"},{\"id\":2,\"from_id\":2,\"to_to\":1,\"rate\":0.1,\"created_at\":\"0001-01-01T00:00:00Z\"}]}\n",
		},
		{
			name: "error request body not json",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.ConversionInsertResponsePayload{}, nil)
				return *service
			}(),
			statusCodeExpected: http.StatusUnprocessableEntity,
			wantErr:            true,
			reqBody:            strings.NewReader(``),
			responseBody:       "{\"message\":\"code=400, message=Request body can't be empty\"}\n",
		},
		{
			name: "error sql",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.ConversionInsertResponsePayload{}, domain.ErrSqlError)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			wantErr:            true,
			reqBody:            strings.NewReader(`{"from_id": 1,"to_id": 2,"rate": 29}`),
			responseBody:       "{\"message\":\"sql error\"}\n",
		},
		{
			name: "error request body empty field",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				service.On("Create", mock.Anything, mock.Anything).Return(domain.ConversionInsertResponsePayload{}, nil)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			wantErr:            true,
			reqBody:            strings.NewReader(`{"to_id": 2,"rate": 29}`),
			responseBody:       "{\"message\":\"Key: 'ConversionInsertRequestPayload.FromID' Error:Field validation for 'FromID' failed on the 'required' tag\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewConversionHandler(e, &tt.mockConversionService)
			req := httptest.NewRequest(http.MethodPost, "/conversion", tt.reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, tt.statusCodeExpected, rec.Code)
			require.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}

func TestConversionHandler_GetRate(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name                  string
		mockConversionService mocks.ConversionService
		path                  string
		statusCodeExpected    int
		wantErr               bool
		responseBody          string
	}{
		{
			name: "success",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				service.On("GetConvertedCurrency", mock.Anything, mock.Anything).Return(domain.ConvertedCurrencyResponse{Result: 100}, nil)
				return *service
			}(),
			statusCodeExpected: http.StatusOK,
			path:               "1/2/10",
			wantErr:            false,
			responseBody:       "{\"result\":100}\n",
		},
		{
			name: "error from_id not exists",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				service.On("GetConvertedCurrency", mock.Anything, mock.Anything).Return(domain.ConvertedCurrencyResponse{}, domain.ErrCurrencyFromIDNotExists)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			path:               "1/2/10",
			wantErr:            true,
			responseBody:       "{\"message\":\"currency from_id not  exists\"}\n",
		},
		{
			name: "error from_id not valid",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			path:               "k/2/10",
			wantErr:            true,
			responseBody:       "{\"message\":\"strconv.ParseInt: parsing \\\"k\\\": invalid syntax\"}\n",
		},
		{
			name: "error from_id not valid",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			path:               "1/c/10",
			wantErr:            true,
			responseBody:       "{\"message\":\"strconv.ParseInt: parsing \\\"c\\\": invalid syntax\"}\n",
		},
		{
			name: "error from_id not valid",
			mockConversionService: func() mocks.ConversionService {
				service := new(mocks.ConversionService)
				return *service
			}(),
			statusCodeExpected: http.StatusBadRequest,
			path:               "1/2/zz",
			wantErr:            true,
			responseBody:       "{\"message\":\"strconv.ParseFloat: parsing \\\"zz\\\": invalid syntax\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewConversionHandler(e, &tt.mockConversionService)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/conversion/rate/%s", tt.path), nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, tt.statusCodeExpected, rec.Code)
			require.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}
