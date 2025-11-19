package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Page struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePageRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content"`
}

type UpdatePageRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// GetAllPages retrieves all pages from the database
func GetAllPages(db *sql.DB) ([]Page, error) {
	rows, err := db.Query("SELECT id, name, content, created_at, updated_at FROM pages ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []Page
	for rows.Next() {
		var page Page
		if err := rows.Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

// GetPageByID retrieves a page by its ID
func GetPageByID(db *sql.DB, id string) (*Page, error) {
	var page Page
	err := db.QueryRow(
		"SELECT id, name, content, created_at, updated_at FROM pages WHERE id = $1",
		id,
	).Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// GetPageByName retrieves a page by its name
func GetPageByName(db *sql.DB, name string) (*Page, error) {
	var page Page
	err := db.QueryRow(
		"SELECT id, name, content, created_at, updated_at FROM pages WHERE name = $1",
		name,
	).Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// CreatePage creates a new page
func CreatePage(db *sql.DB, req CreatePageRequest) (*Page, error) {
	id := uuid.New().String()
	now := time.Now()

	var page Page
	err := db.QueryRow(
		"INSERT INTO pages (id, name, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, content, created_at, updated_at",
		id, req.Name, req.Content, now, now,
	).Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &page, nil
}

// UpdatePage updates an existing page
func UpdatePage(db *sql.DB, id string, req UpdatePageRequest) (*Page, error) {
	now := time.Now()

	// Build dynamic update query
	query := "UPDATE pages SET updated_at = $1"
	args := []interface{}{now}
	argCount := 2

	if req.Name != "" {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, req.Name)
		argCount++
	}

	if req.Content != "" {
		query += fmt.Sprintf(", content = $%d", argCount)
		args = append(args, req.Content)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, content, created_at, updated_at", argCount)
	args = append(args, id)

	var page Page
	err := db.QueryRow(query, args...).Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// DeletePage deletes a page by ID
func DeletePage(db *sql.DB, id string) error {
	result, err := db.Exec("DELETE FROM pages WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
