package api

import (
	"context"
	"fmt"
	sv "github.com/cheojeg/top_phrases/core/services"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	"github.com/cheojeg/top_phrases/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	DraftPhraseState     = "draft"
	PublishedPhraseState = "published"
	ArchivedPhraseState  = "archived"
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
		State:     DraftPhraseState,
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
	ctx.Header("HX-Redirect", "/quotes")
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

	ctx.Header("HX-Redirect", "/quotes")
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
	ctx.Header("HX-Redirect", "/quotes")
	ctx.JSON(http.StatusOK, phrase)
}

func publishPhrase(service sv.Service, config util.Config) {
	ctx := context.Background()
	if config.PostEnabled {
		phrase, err := service.GetPhraseToPublish(ctx, 15)
		if err != nil {
			log.Println("cannot get phrase to publish:", err)
			return
		}

		client, err := newOAuth1Client(config.GotwiApiKey, config.GotwiApiKeySecret, config.TwAccessToken, config.TwAccessSecret)
		if err != nil {
			log.Println("Error creating OAuth1Client:", err)
			return
		}

		log.Println(phrase)
		tweetId, err := xPost(client, phrase)
		if err != nil {
			log.Println("Error posting:", err)
			log.Println(os.Stderr, err)
			return
		}

		log.Println("XPost id", tweetId)
	} else {
		log.Println("Posting is disabled")
	}

}

func newOAuth1Client(apiKey, apiSecret, accessToken, accessSecret string) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               apiKey,
		APIKeySecret:         apiSecret,
		OAuthToken:           accessToken,
		OAuthTokenSecret:     accessSecret,
	}

	return gotwi.NewClient(in)
}

func xPost(c *gotwi.Client, text string) (string, error) {
	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}

func (server *Server) publishQuoteOfTheDay(ctx *gin.Context) {

	if server.config.PostEnabled {
		phrase, err := server.service.GetPhraseToPublish(ctx, 15)
		if err != nil {
			log.Println("cannot get phrase to publish:", err)
			ctx.Header("HX-Redirect", "/check_quote_of_the_day")
			err := fmt.Errorf("cannot get phrase to publish:", err)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		client, err := newOAuth1Client(server.config.GotwiApiKey, server.config.GotwiApiKeySecret, server.config.TwAccessToken, server.config.TwAccessSecret)
		if err != nil {
			log.Println("Error creating OAuth1Client:", err)
			return
		}

		log.Println(phrase)
		tweetId, err := xPost(client, phrase)
		if err != nil {
			log.Println("Error posting:", err)
			log.Println(os.Stderr, err)
			return
		}

		log.Println("tweet id", tweetId)

		ctx.Header("HX-Redirect", "/check_quote_of_the_day")
	} else {
		log.Println("Posting is disabled")
		ctx.Header("HX-Redirect", "/check_quote_of_the_day")
		err := fmt.Errorf("Posting is disabled")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
