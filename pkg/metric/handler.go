package metric

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
}

func (h *Handler) Register(router *gin.Engine) {
	router.GET(URL, h.Heartbeat)
}

// Heartbeat
// @Summary Heartbeat metric
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/heartbeat [get]
func (h *Handler) Heartbeat(g *gin.Context) {
	g.Writer.WriteHeader(http.StatusNoContent)
}
