package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

type Quote struct {
	Text   string
	Author string
}

func (server *Server) quotes(ctx *gin.Context) {
	quotes := []Quote{
		{Text: "The only limit to our realization of tomorrow is our doubts of today.", Author: "Franklin D. Roosevelt"},
		{Text: "The purpose of our lives is to be happy.", Author: "Dalai Lama"},
		{Text: "Life is what happens when you're busy making other plans.", Author: "John Lennon"},
	}
	ctx.HTML(http.StatusOK, "quotes.html", gin.H{
		"title":  "Quotes",
		"Quotes": quotes,
	})
}
