package domain

import (
	"context"
	"time"
)

type Currency struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CurrencyResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CurrenciesResponsePayload struct {
	List []CurrencyResponse `json:"list"`
}

type CurrencyRequestPayload struct {
	Name string `json:"name" validate:"required"`
}

type CurrencyResponsePayload struct {
	Name string `json:"name" `
}

type CurrencyRepository interface {
	GetAll(ctx context.Context) (res []Currency, err error)
	Create(ctx context.Context, comment *Currency) (err error)
	CheckIsExistByName(ctx context.Context, orgName string) (isExist bool, err error)
	CheckIsExistByID(ctx context.Context, id int64) (isAlreadyExist bool, err error)
}

type CurrencyService interface {
	GetAll(ctx context.Context) (response CurrenciesResponsePayload, err error)
	Create(ctx context.Context, request CurrencyRequestPayload) (currency Currency, err error)
}
