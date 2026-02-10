package frontend

import (
	"net/http"
	"zugzwang/internal/ui/layouts"
	"zugzwang/internal/ui/pages"

	"github.com/gin-gonic/gin"
)

func HandleIndex(c *gin.Context) {

	// Set HTTP header
	c.Header("Content-Type", "text/html; charset=utf-8")

	// Build page
	page := layouts.MainLayout(pages.Index(), "Zugzwang")

	// Render the page
	err := page.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to render index page.")
	}

}
