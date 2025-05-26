package repository

import (
	"BlockchainCurrency/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type QuotesRepository struct {
	db *sql.DB
}

func NewQuotesRepository(db *sql.DB) *QuotesRepository {
	return &QuotesRepository{db: db}
}

func (q *QuotesRepository) GetAllQuotes(ctx context.Context) ([]domain.Quotes, error) {
	query := `SELECT quotation,price_24h,volume_24h,last_trade_price FROM quotation_s`

	var quotations []domain.Quotes

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query quotations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var qt domain.Quotes
		err := rows.Scan(

			&qt.Symbol,
			&qt.Price24h,
			&qt.Volume24h,
			&qt.LastTradePrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quotation row: %w", err)
		}
		quotations = append(quotations, qt)
	}

	return quotations, nil
}

func (q *QuotesRepository) SaveQuotes(ctx context.Context, quotes []domain.Quotes) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
        INSERT INTO quotation_s (quotation, price_24h, volume_24h, last_trade_price)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (quotation) DO UPDATE SET
            price_24h = EXCLUDED.price_24h,
            volume_24h = EXCLUDED.volume_24h,
            last_trade_price = EXCLUDED.last_trade_price
    `

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, quote := range quotes {
		if _, err := stmt.ExecContext(ctx,
			quote.Symbol,
			quote.Price24h,
			quote.Volume24h,
			quote.LastTradePrice,
		); err != nil {
			return fmt.Errorf("failed to save %s: %w", quote.Symbol, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
