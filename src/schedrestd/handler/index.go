package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler ...
type Handler struct {

}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}
