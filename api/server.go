package api

import (
	"fmt"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	"github.com/cheojeg/top_phrases/token"
	"github.com/gin-gonic/gin"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	//tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		layoutCopy = append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), layoutCopy...)
	}
	return r
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/users/login_web", server.loginUserWeb)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/phrase", server.createPhrase)
	authRoutes.PUT("/phrase", server.updatePhrase)
	authRoutes.PUT("/phrase_state", server.updatePhraseState)
	authRoutes.GET("/quotes", server.quotes)
	authRoutes.GET("/create_quote", server.createQuote)
	authRoutes.POST("/create_quote", server.createQuoteWeb)
	authRoutes.GET("/edit_quote/:id", server.editQuoteWeb)
	authRoutes.GET("/update_state_quote/:id", server.updateStateQuoteWeb)

	router.HTMLRender = loadTemplates("./templates")
	router.GET("/", server.index)

	server.router = router
}

func (server *Server) Start(address string) error {
	//address = "0.0.0.0:8080"
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
