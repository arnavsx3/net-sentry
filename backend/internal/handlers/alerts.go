package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/repository"
)

func GetCurrentAlerts(repo *repository.AlertRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 50
		if raw := c.Query("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid limit",
				})
				return
			}
			limit = parsed
		}

		items, err := repo.GetCurrentAlerts(c.Request.Context(), limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to fetch current alerts",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count":  len(items),
			"alerts": items,
		})
	}
}