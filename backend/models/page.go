package models

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Page struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PageLink struct {
	ID         int       `json:"id"`
	FromPageID int       `json:"from_page_id"`
	ToPageID   int       `json:"to_page_id"`
	LinkText   string    `json:"link_text"`
	CreatedAt  time.Time `json:"created_at"`
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
func GetPageByID(db *sql.DB, id int) (*Page, error) {
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
	now := time.Now()

	var page Page
	err := db.QueryRow(
		"INSERT INTO pages (name, content, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name, content, created_at, updated_at",
		req.Name, req.Content, now, now,
	).Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &page, nil
}

// UpdatePage updates an existing page
func UpdatePage(db *sql.DB, id int, req UpdatePageRequest) (*Page, error) {
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
func DeletePage(db *sql.DB, id int) error {
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

// ExtractWikiLinks extracts page names from [[Page Name]] syntax in content
func ExtractWikiLinks(content string) []string {
	wikiLinkRegex := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	matches := wikiLinkRegex.FindAllStringSubmatch(content, -1)
	
	links := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	
	for _, match := range matches {
		if len(match) > 1 {
			pageName := strings.TrimSpace(match[1])
			if pageName != "" && !seen[pageName] {
				links = append(links, pageName)
				seen[pageName] = true
			}
		}
	}
	
	return links
}

// UpdatePageLinks updates the page_links table for a given page
func UpdatePageLinks(db *sql.DB, pageID int, content string) error {
	// Extract links from content
	linkNames := ExtractWikiLinks(content)
	
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete existing links for this page
	_, err = tx.Exec("DELETE FROM page_links WHERE from_page_id = $1", pageID)
	if err != nil {
		return err
	}
	
	// Insert new links
	for _, linkName := range linkNames {
		// Find the target page by name
		var toPageID int
		err := tx.QueryRow("SELECT id FROM pages WHERE name = $1", linkName).Scan(&toPageID)
		if err == sql.ErrNoRows {
			// Page doesn't exist yet, skip this link
			continue
		}
		if err != nil {
			return err
		}
		
		// Insert link (ignore if duplicate due to UNIQUE constraint)
		_, err = tx.Exec(
			"INSERT INTO page_links (from_page_id, to_page_id, link_text) VALUES ($1, $2, $3) ON CONFLICT (from_page_id, to_page_id) DO NOTHING",
			pageID, toPageID, linkName,
		)
		if err != nil {
			return err
		}
	}
	
	// Commit transaction
	return tx.Commit()
}

// GetBacklinks retrieves all pages that link to the given page
func GetBacklinks(db *sql.DB, pageID int) ([]Page, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, p.content, p.created_at, p.updated_at
		FROM pages p
		INNER JOIN page_links pl ON p.id = pl.from_page_id
		WHERE pl.to_page_id = $1
		ORDER BY p.name
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	pages := make([]Page, 0)
	for rows.Next() {
		var page Page
		if err := rows.Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	
	return pages, nil
}

// GetForwardLinks retrieves all pages that the given page links to
func GetForwardLinks(db *sql.DB, pageID int) ([]Page, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, p.content, p.created_at, p.updated_at
		FROM pages p
		INNER JOIN page_links pl ON p.id = pl.to_page_id
		WHERE pl.from_page_id = $1
		ORDER BY p.name
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	pages := make([]Page, 0)
	for rows.Next() {
		var page Page
		if err := rows.Scan(&page.ID, &page.Name, &page.Content, &page.CreatedAt, &page.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	
	return pages, nil
}

// UpdateLinksToNewPage updates all existing pages that reference the newly created page
// This ensures that backlinks are created when a page is created from a missing link
func UpdateLinksToNewPage(db *sql.DB, newPageName string, newPageID int) error {
	// Get all pages and check which ones reference the new page
	allPages, err := GetAllPages(db)
	if err != nil {
		return err
	}
	
	// For each page, check if it contains a wiki link to the new page
	for _, page := range allPages {
		// Skip if this is the new page itself
		if page.ID == newPageID {
			continue
		}
		
		// Extract wiki links from this page's content
		linkNames := ExtractWikiLinks(page.Content)
		
		// Check if this page references the new page
		referencesNewPage := false
		for _, linkName := range linkNames {
			if linkName == newPageName {
				referencesNewPage = true
				break
			}
		}
		
		// If this page references the new page, update its links
		// This will now create the link since the target page exists
		if referencesNewPage {
			if err := UpdatePageLinks(db, page.ID, page.Content); err != nil {
				// Log error but continue with other pages
				// We don't want to fail the entire operation if one page update fails
				continue
			}
		}
	}
	
	return nil
}
