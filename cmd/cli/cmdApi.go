package cli

import (
	"database/sql"
	"github.com/cheojeg/top_phrases/api"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
			conn, err := sql.Open(config.DBDriver, config.DBSource)

			// Run db migration
			runDbMigration(config.MigrationURL, config.DBSource)

			store := db.NewStore(conn)
			runGinServer(*config, store)

			return nil
		},
	}
	return cmd
}
