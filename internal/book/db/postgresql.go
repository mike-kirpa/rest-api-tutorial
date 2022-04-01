package db

import (
	"context"
	"fmt"
	"restapi-lesson/internal/book"
	"restapi-lesson/pkg/client/postgresql"
	"restapi-lesson/pkg/logging"

	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) book.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) FindAll(ctx context.Context) (u []book.Book, err error) {
	q := `SELECT id, name, age FROM public.book`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	books := make([]book.Book, 0)

	for rows.Next() {
		var bk Book

		err = rows.Scan(&bk.ID, &bk.Name, &bk.Age)
		if err != nil {
			return nil, err
		}
		/*		sq := `
					SELECT
						a.id, a.name
					FROM book_authors
					JOIN public.author a on a.id = book_authors.author_id
					WHERE book_id = $1;
					`
				authorsRows, err := r.client.Query(ctx, sq, bk.ID)
				if err != nil {
					return nil, err
				}

				authors := make([]author.Author, 0)

				for authorsRows.Next() {
					var ath author.Author
					err = authorsRows.Scan(&ath.ID, &ath.Name)
					if err != nil {
						return nil, err
					}
					authors = append(authors, ath)
				}

				bk.Authors = authors
		*/
		books = append(books, bk.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
