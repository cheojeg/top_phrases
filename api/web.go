package api

import (
	db "github.com/cheojeg/top_phrases/db/sqlc"
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

func parseQuotes(phrases []db.Phrase) []Quote {
	quotes := make([]Quote, len(phrases))
	for i, phrase := range phrases {
		quotes[i] = Quote{
			Text:   phrase.Phrase,
			Author: phrase.Author,
		}
	}
	return quotes
}

func (server *Server) quotes(ctx *gin.Context) {
	quotesQuery, err := server.store.ListPhrases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	quotes := parseQuotes(quotesQuery)
	ctx.HTML(http.StatusOK, "quotes.html", gin.H{
		"title":  "Quotes",
		"Quotes": quotes,
	})
}
