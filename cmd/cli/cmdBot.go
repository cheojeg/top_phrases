package cli

import (
	"context"
	"database/sql"
	sv "github.com/cheojeg/top_phrases/core/services"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func newCmdBot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bot",
		Short: "Starts the API server",
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

			botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
			if botToken == "" {
				log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
			}

			chatID := os.Getenv("CHAT_ID")
			if chatID == "" {
				log.Fatal("CHAT_ID environment variable is not set")
			}

			bot, err := tgbotapi.NewBotAPI(botToken)
			if err != nil {
				log.Panic(err)
			}

			bot.Debug = true

			log.Printf("Authorized on account %s", bot.Self.UserName)
			for {
				// Select a random message
				ctx := context.Background()
				phrase, err := service.GetPhraseToPublish(ctx)
				//phrase, err := store.GetPhraseToPublish(ctx)
				if err != nil {
					log.Fatal("cannot get phrase to publish:", err)
					return nil
				}

				msg := tgbotapi.NewMessageToChannel(chatID, phrase)
				msg.ParseMode = "MarkdownV2"
				bot.Send(msg)

				// Sleep for 60 seconds
				time.Sleep(60 * time.Second)
			}
			return nil
		},
	}
	return cmd
}
