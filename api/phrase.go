package api

import (
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

const (
	draftPhraseState     = "draft"
	publishedPhraseState = "published"
)

type createPhraseRequest struct {
	Phrase string `json:"phrase" binding:"required"`
	Author string `json:"author"`
}

func (server *Server) createPhrase(ctx *gin.Context) {
	var req createPhraseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreatePhraseParams{
		Owner:     authPayload.Username,
		State:     draftPhraseState,
		Phrase:    req.Phrase,
		Author:    req.Author,
		CreatedAt: time.Now(),
	}

	phrase, err := server.store.CreatePhrase(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, phrase)
}

type updatePhraseStateRequest struct {
	ID    int64  `json:"id" binding:"required"`
	State string `json:"state" binding:"required"`
}

func (server *Server) updatePhraseState(ctx *gin.Context) {
	var req updatePhraseStateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePhraseStateParams{
		ID:    req.ID,
		State: req.State,
	}
	phrase, err := server.store.UpdatePhraseState(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, phrase)
}

type updatePhraseRequest struct {
	ID     int64  `json:"id" binding:"required"`
	Phrase string `json:"phrase" binding:"required"`
	Author string `json:"author"`
}

func (server *Server) updatePhrase(ctx *gin.Context) {
	var req updatePhraseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePhraseParams{
		ID:     req.ID,
		Phrase: req.Phrase,
		Author: req.Author,
	}

	phrase, err := server.store.UpdatePhrase(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, phrase)
}
