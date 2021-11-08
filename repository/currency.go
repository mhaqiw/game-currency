package repository

import (
	"context"
	"database/sql"
	"github.com/labstack/gommon/log"
	"github.com/mhaqiw/game-currency/domain"
	"time"
)

type currencyRepository struct {
	Conn *sql.DB
}

func NewCurrencyRepository(Conn *sql.DB) domain.CurrencyRepository {
	return &currencyRepository{Conn}
}

func (p currencyRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Currency, error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	result := make([]domain.Currency, 0)
	for rows.Next() {
		t := domain.Currency{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (p *currencyRepository) GetAll(ctx context.Context) (res []domain.Currency, err error) {
	query := `SELECT id, name, created_at FROM currency order by id`
	res, err = p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}

func (p *currencyRepository) Create(ctx context.Context, curr *domain.Currency) (err error) {
	curr.CreatedAt = time.Now()
	query := `INSERT INTO currency( name, created_at) VALUES( $1, $2 ) RETURNING id`
	err = p.Conn.QueryRowContext(ctx, query, curr.Name, curr.CreatedAt).Scan(&curr.ID)
	if err != nil {
		return
	}

	return
}

func (p *currencyRepository) CheckIsExistByName(ctx context.Context, orgName string) (isAlreadyExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM currency WHERE name= $1)`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	row := stmt.QueryRowContext(ctx, orgName)
	if err != nil {
		return false, err
	}
	err = row.Scan(
		&isAlreadyExist,
	)
	return
}

func (p *currencyRepository) CheckIsExistByID(ctx context.Context, id int64) (isAlreadyExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM currency WHERE id= $1)`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	row := stmt.QueryRowContext(ctx, id)
	if err != nil {
		return false, err
	}
	err = row.Scan(
		&isAlreadyExist,
	)
	return
}