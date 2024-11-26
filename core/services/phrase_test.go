package services

import (
	"context"
	"database/sql"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"github.com/cheojeg/top_phrases/db/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestCreatePhrase(t *testing.T) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := db.NewStore(conn)
	service := NewService(store)

	phrase, err := service.GetPhraseToPublish(context.Background(), 4)
	// if err != nil {
	//	log.Fatal("cannot get a quote to publish:", err)
	// }
	require.NotEmpty(t, phrase)
}
