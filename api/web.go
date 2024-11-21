package api

import (
	"database/sql"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (server *Server) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

type Quote struct {
	ID     int64
	Text   string
	Author string
	State  string
}

func parseQuotes(phrases []db.Phrase) []Quote {
	quotes := make([]Quote, len(phrases))
	for i, phrase := range phrases {
		quotes[i] = Quote{
			ID:     phrase.ID,
			Text:   phrase.Phrase,
			Author: phrase.Author,
			State:  phrase.State,
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
	ctx.Header("HX-Redirect", "/create_quote")
	ctx.HTML(http.StatusOK, "create_quote.html", gin.H{
		"title": "Crear Frase",
	})
}

type createQuoteRequest struct {
	Quote  string `form:"quote" binding:"required"`
	Author string `form:"author"`
}

func (server *Server) createQuoteWeb(ctx *gin.Context) {
	var req createQuoteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreatePhraseParams{
		Owner:     authPayload.Username,
		State:     draftPhraseState,
		Phrase:    req.Quote,
		Author:    req.Author,
		CreatedAt: time.Now(),
	}

	_, err := server.store.CreatePhrase(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Header("HX-Redirect", "/quotes")
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
	//quotes := parseQuotes(quotesQuery)
	ctx.HTML(http.StatusOK, "edit_quote.html", gin.H{
		"title": "Edit Quote",
		"Quote": quote,
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
	//quotes := parseQuotes(quotesQuery)
	ctx.HTML(http.StatusOK, "state_quote.html", gin.H{
		"title": "Edit State Quote",
		"Quote": quote,
	})
}
