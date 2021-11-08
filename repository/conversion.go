package repository

import (
	"context"
	"database/sql"
	"github.com/labstack/gommon/log"
	"github.com/mhaqiw/game-currency/domain"
	"time"
)

type conversionRepository struct {
	Conn *sql.DB
}

func NewConversionRepository(Conn *sql.DB) domain.ConversionRepository {
	return &conversionRepository{Conn}
}


func (p *conversionRepository) Post(ctx context.Context, conversion domain.Conversion) (fist int64, second int64, createdAt time.Time, err error) {
	query := `INSERT INTO conversion(from_id, to_id,rate ,created_at) VALUES ($1 , $2 , $3, $4) RETURNING id`
	createdAt = time.Now()
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	err = p.Conn.QueryRowContext(ctx, query, conversion.FromID, conversion.ToID, conversion.Rate, createdAt).Scan(&fist)
	if err != nil {
		return
	}

	err = p.Conn.QueryRowContext(ctx, query, conversion.ToID, conversion.FromID, 1/conversion.Rate, createdAt).Scan(&second)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func (p *conversionRepository) GetFromToID(ctx context.Context, fromID int64, toID int64) (res domain.Conversion, err error) {
	query := `SELECT id, from_id, to_id, rate, created_at FROM conversion WHERE from_id= $1 and to_id= $2`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	row := stmt.QueryRowContext(ctx, fromID, toID)
	if err != nil {
		return
	}
	err = row.Scan(
		&res.ID, &res.FromID, &res.ToID, &res.Rate, &res.CreatedAt,
	)
	return
}
