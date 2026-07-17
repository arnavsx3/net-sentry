package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/repository"
)

func GetLatestTraceroute(repo *repository.TracerouteRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Param("host")
		if host == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "target host is required",
			})
			return
		}

		item, err := repo.GetLatestTraceroute(c.Request.Context(), host)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to fetch latest traceroute",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}