package repository

import (
	"database/sql"
	"github.com/youngabreu1/knowledge-base.git/internal/domain"
)

type PageRepository interface {
	FindAll() ([]domain.Page, error)
}

type mysqlPageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) PageRepository{
	return &mysqlPageRepository {
		db: db,
	}
}
func (r *mysqlPageRepository) FindAll() ([]domain.Page, error) {
	rows, err := r.db.Query("SELECT text FROM pages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []domain.Page
	
	for rows.Next() {
		var page domain.Page
		if err := rows.Scan(&page.Text); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}