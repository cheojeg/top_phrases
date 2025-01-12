package cli

import (
	"context"
	"database/sql"
	"fmt"
	sv "github.com/cheojeg/top_phrases/core/services"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

func newCmdBot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bot",
		Short: "Starts the bot server",
		RunE: func(cmd *cobra.Command, args []string) error {

			config, err := util.LoadConfig(".")
			if err != nil {
				log.Fatal("cannot load config:", err)
			}
			conn, err := sql.Open(config.DBDriver, config.DBSource)
			if err != nil {
				log.Fatal("cannot connect to database:", err)
			}
			store := db.NewStore(conn)
			service := sv.NewService(store)

			if config.TwAccessToken == "" {
				log.Fatal("TW_ACCESS_TOKEN environment variable is not set")
			}

			if config.TwAccessSecret == "" {
				log.Fatal("TW_ACCESS_SECRET environment variable is not set")
			}

			if config.TwAccessToken == "" || config.TwAccessSecret == "" {
				log.Fatal("Please set the TW_ACCESS_TOKEN and TW_ACCESS_SECRET environment variables.")
				os.Exit(1)
			}

			for {
				// Select a random message
				ctx := context.Background()
				phrase, err := service.GetPhraseToPublish(ctx, 15)
				if err != nil {
					log.Println("cannot get phrase to publish:", err)
					time.Sleep(30 * time.Minute)
					//time.Sleep(24 * time.Hour)
					continue
					//return nil
				}

				client, err := newOAuth1Client(config.TwAccessToken, config.TwAccessSecret)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				log.Println(phrase)
				tweetId, err := tweet(client, phrase)
				if err != nil {
					log.Println(os.Stderr, err)
					//os.Exit(2)
				}

				log.Println("tweet id", tweetId)
				// Sleep for 24 hours
				time.Sleep(24 * time.Hour)
			}
			return nil
		},
	}
	return cmd
}

func newOAuth1Client(accessToken, accessSecret string) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           accessToken,
		OAuthTokenSecret:     accessSecret,
	}

	return gotwi.NewClient(in)
}

func tweet(c *gotwi.Client, text string) (string, error) {
	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}
