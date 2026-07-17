package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/arnavsx3/net-sentry/backend/internal/realtime"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWebSocket(hub *realtime.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to upgrade websocket connection",
			})
			return
		}

		hub.Register(conn)

		_ = conn.WriteJSON(gin.H{
			"type":      "connected",
			"timestamp": time.Now().UTC(),
			"payload": gin.H{
				"message": "websocket connected",
			},
		})

		go func() {
			defer hub.Unregister(conn)

			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	}
}