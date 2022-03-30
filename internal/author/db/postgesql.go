package author

import (
	"context"
	"fmt"
	"restapi-lesson/internal/author"
	"restapi-lesson/pkg/client/postgresql"
	"restapi-lesson/pkg/logging"
	"strings"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

// Create implements author.Repository
func (r *repository) Create(ctx context.Context, author *author.Author) error {
	q := `INSERT INTO author (name) VALUES ($1) RETURNING id`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, author.Name).Scan(&author.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code, pgerr.SQLState())
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

// Delete implements author.Repository
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("dfd")
}

// FindAll implements author.Repository
func (r *repository) FindAll(ctx context.Context) (u []author.Author, err error) {
	q := `
		SELECT id, name FROM public.author;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	authors := make([]author.Author, 0)

	for rows.Next() {
		var ath author.Author

		err = rows.Scan(&ath.ID, &ath.Name)
		if err != nil {
			return nil, err
		}

		authors = append(authors, ath)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

// FindOne implements author.Repository
func (r *repository) FindOne(ctx context.Context, id string) (author.Author, error) {
	q := `SELECT id, name FROM public.author WHERE id = $1`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	var auth author.Author
	err := r.client.QueryRow(ctx, q, id).Scan(&auth.ID, &auth.Name)
	if err != nil {
		return author.Author{}, err
	}
	return auth, nil
}

// Update implements author.Repository
func (r *repository) Update(ctx context.Context, author author.Author) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) author.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
