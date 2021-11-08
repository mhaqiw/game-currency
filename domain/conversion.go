package domain

import (
	"context"
	"time"
)

type Conversion struct {
	ID        int64     `json:"id"`
	FromID    int64     `json:"from_id"`
	ToID      int64     `json:"to_to"`
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversionInsertRequestPayload struct {
	FromID int64   `json:"from_id" validate:"required"`
	ToID   int64   `json:"to_id" validate:"required"`
	Rate   float64 `json:"rate" validate:"required"`
}

type ConversionInsertResponsePayload struct {
	List []Conversion `json:"list"`
}

type ConvertedCurrencyResponse struct {
	Result float64 `json:"result"`
}

type ConvertedCurrencyRequest struct {
	FromID int64   `json:"from_id"`
	ToID   int64   `json:"to_id"`
	Amount float64 `json:"amount"`
}

type ConversionRepository interface {
	GetFromToID(ctx context.Context, fromID int64, toID int64) (res Conversion, err error)
	Post(ctx context.Context, conversion Conversion) (fist int64, second int64, createdAt time.Time, err error)
}

type ConversionService interface {
	Create(ctx context.Context, request ConversionInsertRequestPayload) (response ConversionInsertResponsePayload, err error)
	GetConvertedCurrency(ctx context.Context, request ConvertedCurrencyRequest) (response ConvertedCurrencyResponse, err error)
}
