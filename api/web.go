package api

import (
	"database/sql"
	"fmt"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (server *Server) index(ctx *gin.Context) {
	log.Println("Hello this is index")
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

type Quote struct {
	ID          int64
	Text        string
	Author      string
	State       string
	PublishedAt string
}

func sqlNullTimeToString(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

func parseQuotes(phrases []db.Phrase) []Quote {
	quotes := make([]Quote, len(phrases))
	for i, phrase := range phrases {
		quotes[i] = Quote{
			ID:          phrase.ID,
			Text:        phrase.Phrase,
			Author:      phrase.Author,
			State:       phrase.State,
			PublishedAt: sqlNullTimeToString(phrase.PublishedAt),
		}
	}
	return quotes
}

func isValidState(state string) bool {
	return state == PublishedPhraseState || state == ArchivedPhraseState || state == DraftPhraseState
}

func (server *Server) quotes(ctx *gin.Context) {

	state := ctx.Query("state")
	fmt.Println(state)
	quotesQuery := []db.Phrase{}
	var err error
	if isValidState(state) {
		quotesQuery, err = server.store.ListPhrasesByState(ctx, state)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		quotesQuery, err = server.store.ListPhrases(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	quotes := parseQuotes(quotesQuery)
	countDraft, err := server.store.CountDraftPhrases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.HTML(http.StatusOK, "quotes.html", gin.H{
		"title":   "Quotes",
		"Pending": countDraft,
		"Quotes":  quotes,
	})
}

func (server *Server) createQuote(ctx *gin.Context) {
	countDraft, err := server.store.CountDraftPhrases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Header("HX-Redirect", "/create_quote")
	ctx.HTML(http.StatusOK, "create_quote.html", gin.H{
		"title":   "Crear Frase",
		"Pending": countDraft,
	})
}

func (server *Server) editQuoteWeb(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	phrase, err := server.store.GetPhraseByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	quote := Quote{
		ID:     phrase.ID,
		Text:   phrase.Phrase,
		Author: phrase.Author,
		State:  phrase.State,
	}

	countDraft, err := server.store.CountDraftPhrases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//quotes := parseQuotes(quotesQuery)
	ctx.HTML(http.StatusOK, "edit_quote.html", gin.H{
		"title":   "Edit Quote",
		"Quote":   quote,
		"Pending": countDraft,
	})
}

func (server *Server) updateStateQuoteWeb(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	phrase, err := server.store.GetPhraseByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	quote := Quote{
		ID:     phrase.ID,
		Text:   phrase.Phrase,
		Author: phrase.Author,
		State:  phrase.State,
	}

	countDraft, err := server.store.CountDraftPhrases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//quotes := parseQuotes(quotesQuery)
	ctx.HTML(http.StatusOK, "state_quote.html", gin.H{
		"title":   "Edit State Quote",
		"Quote":   quote,
		"Pending": countDraft,
	})
}

func (server *Server) inboxQuotes(ctx *gin.Context) {
	quotesQuery, err := server.store.ListPhrasesByState(ctx, DraftPhraseState)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	quotes := parseQuotes(quotesQuery)
	countDraft, err := server.store.CountDraftPhrases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.HTML(http.StatusOK, "inbox_quotes.html", gin.H{
		"title":   "Quotes",
		"Pending": countDraft,
		"Quotes":  quotes,
	})
}
