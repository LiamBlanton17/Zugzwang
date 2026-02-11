package frontend

import (
	"net/http"
	"zugzwang/internal/ui/layouts"
	"zugzwang/internal/ui/pages"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func HandleIndex(c *gin.Context) {
	// Check if it's an HTMX request
	isHTMX := c.GetHeader("HX-Request") == "true"

	// Get component, wrap if not htmx
	var component templ.Component
	if isHTMX {
		component = pages.Index()
	} else {
		component = layouts.MainLayout(pages.Index(), "ZugZwang")
	}

	// Set Header and Render
	c.Header("Content-Type", "text/html; charset=utf-8")
	err := component.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to render.")
	}

}
