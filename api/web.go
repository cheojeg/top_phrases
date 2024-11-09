package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website for top quotes",
	})
}
