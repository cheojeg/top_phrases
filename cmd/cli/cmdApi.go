package cli

import (
	"context"
	"database/sql"
	"github.com/cheojeg/top_phrases/api"
	sv "github.com/cheojeg/top_phrases/core/services"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func runGinServer(config util.Config, store db.Store, service sv.Service) {
	server, err := api.NewServer(config, store, service)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runDbMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create migrate instance:", err)
	}

	if migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrate up:", err)
	}

	log.Println("db Migrated successfully!")
}

func newCmdApi() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Starts the API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := util.LoadConfig(".")
			if err != nil {
				log.Fatal("cannot load config:", err)
			}
			log.Println(config.DBSource)
			conn, err := sql.Open(config.DBDriver, config.DBSource)

			// Run db migration
			runDbMigration(config.MigrationURL, config.DBSource)

			store := db.NewStore(conn)
			service := sv.NewService(store)

			c := cron.New()
			c.AddFunc("@every 24h", func() { publishPhrase(*service, *config) })
			c.Start()

			runGinServer(*config, store, *service)

			return nil
		},
	}
	return cmd
}

func publishPhrase(service sv.Service, config util.Config) {
	ctx := context.Background()
	if config.PostEnabled {
		phrase, err := service.GetPhraseToPublish(ctx, 15)
		if err != nil {
			log.Println("cannot get phrase to publish:", err)
			return
		}

		client, err := NewOAuth1Client(config.GotwiApiKey, config.GotwiApiKeySecret, config.TwAccessToken, config.TwAccessSecret)
		if err != nil {
			log.Println("Error creating OAuth1Client:", err)
			return
		}

		log.Println(phrase)
		tweetId, err := XPost(client, phrase)
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

func NewOAuth1Client(apiKey, apiSecret, accessToken, accessSecret string) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               apiKey,
		APIKeySecret:         apiSecret,
		OAuthToken:           accessToken,
		OAuthTokenSecret:     accessSecret,
	}

	return gotwi.NewClient(in)
}

func XPost(c *gotwi.Client, text string) (string, error) {
	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}
