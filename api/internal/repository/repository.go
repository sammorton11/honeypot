package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sammorton11/honeypot-proxy/internal/database"
	"github.com/sammorton11/honeypot-proxy/internal/models"
)

type AttemptStoreInterface interface {
	GetByID(ctx context.Context, id int) (models.Attempt, error)
	GetAll(ctx context.Context) ([]models.Attempt, error)
	Update(ctx context.Context, id int) error
	DeleteByID(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) error
	Insert(ctx context.Context, attempt models.Attempt) error
}

type AttemptStore struct {
	db *pgxpool.Pool
}

func NewAttemptStore(ctx context.Context) (*AttemptStore, error) {
	db, err := database.NewPool(ctx)
	if err != nil {
		return nil, err
	}
	return &AttemptStore{
		db: db,
	}, nil
}

func (s *AttemptStore) Insert(ctx context.Context, attempt models.Attempt) error {
	query := `INSERT INTO attempts (address, network, message) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(ctx, query, attempt.Address, attempt.Network, attempt.Message)
	if err != nil {
		log.Println("Attempt Store: \nAttempt: ", attempt)
		log.Println(err)
		return err
	}
	return nil
}

func (s *AttemptStore) Update(ctx context.Context,id int) error {
	log.Println("Not implemented")
	return nil
}

func (s *AttemptStore) GetByID(ctx context.Context, id int) (models.Attempt, error) {
	query := `SELECT address, network, message FROM attempts WHERE id = $1`
	row := s.db.QueryRow(ctx, query, id)

	var attempt models.Attempt
	err := row.Scan(&attempt.Address, &attempt.Network, &attempt.Message)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Attempt{}, fmt.Errorf("attempt not found %d\n", id)
		}
		return models.Attempt{}, err
	}
	return attempt, nil
}

func (s *AttemptStore) GetAll(ctx context.Context) ([]models.Attempt, error) {
	// select all attempts query
	query := `SELECT address, network, message FROM attempts`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.Attempt
	for rows.Next() {
		var address, network, message string
		err := rows.Scan(&address, &network, &message)
		if err != nil {
			return nil, err
		}

		attempt := models.Attempt{
			Address: address,
			Network: network,
			Message: message,
		}
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

func (s *AttemptStore) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM attempts`
	_, err := s.db.Exec(ctx, query)
	return err
}

func (s *AttemptStore) DeleteByID(ctx context.Context, id int) error {
		query := `DELETE FROM attempts WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	return err
}

