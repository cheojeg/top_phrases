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

	phrase, err := service.GetPhraseToPublish(context.Background())
	if err != nil {
		log.Fatal("cannot get phrase:", err)
	}
	require.NotEmpty(t, phrase)
	//publishedAt := time.Now()
	//require.WithinDuration(t, phrase.P, publishedAt, time.Second)

	//maker, err := NewJWTMaker(util.RandomString(32))
	//require.NoError(t, err)
	//
	//username := util.RandomOwner()
	////role := util.DepositorRole
	//duration := time.Minute
	//
	//issuedAt := time.Now()
	//expiredAt := issuedAt.Add(duration)
	//
	//token, payload, err := maker.CreateToken(username, "role", duration)
	//require.NoError(t, err)
	//require.NotEmpty(t, token)
	//require.NotEmpty(t, payload)
	//
	//payload, err = maker.VerifyToken(token)
	//require.NoError(t, err)
	//require.NotEmpty(t, token)
	//
	//require.NotZero(t, payload.ID)
	//require.Equal(t, username, payload.Username)
	////require.Equal(t, role, payload.Role)
	//require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	//require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
