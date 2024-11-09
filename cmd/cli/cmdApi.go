package cli

import (
	"database/sql"
	"github.com/cheojeg/top_phrases/api"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"log"
)

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
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
			conn, err := sql.Open(config.DBDriver, config.DBSource)
			store := db.NewStore(conn)
			runGinServer(*config, store)

			return nil
		},
	}
	return cmd
}
