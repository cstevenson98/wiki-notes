package handlers

import (
	"database/sql"
	"net/http"

	"github.com/conor/wiki-notes-backend/models"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	DB *sql.DB
}

func NewPageHandler(db *sql.DB) *PageHandler {
	return &PageHandler{DB: db}
}

// GetAllPages handles GET /api/pages
func (h *PageHandler) GetAllPages(c *gin.Context) {
	pages, err := models.GetAllPages(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetPageByID handles GET /api/page/:id
func (h *PageHandler) GetPageByID(c *gin.Context) {
	id := c.Param("id")

	page, err := models.GetPageByID(h.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if page == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// GetPageByName handles GET /api/page/by-name/:name
func (h *PageHandler) GetPageByName(c *gin.Context) {
	name := c.Param("name")

	page, err := models.GetPageByName(h.DB, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if page == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// CreatePage handles POST /api/page
func (h *PageHandler) CreatePage(c *gin.Context) {
	var req models.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := models.CreatePage(h.DB, req)
	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"pages_name_key\"" {
			c.JSON(http.StatusConflict, gin.H{"error": "Page with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// UpdatePage handles PATCH /api/page/:id
func (h *PageHandler) UpdatePage(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := models.UpdatePage(h.DB, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if page == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// DeletePage handles DELETE /api/page/:id
func (h *PageHandler) DeletePage(c *gin.Context) {
	id := c.Param("id")

	err := models.DeletePage(h.DB, id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
