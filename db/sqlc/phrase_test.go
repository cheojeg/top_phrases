package db

import (
	"context"
	"database/sql"
	"github.com/cheojeg/top_phrases/db/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

const (
	draftPhraseState     = "draft"
	publishedPhraseState = "published"
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
	store := NewStore(conn)

	hashedPassword, err := util.HashPassword(util.RandomString(12))
	if err != nil {
		log.Fatal("cannot create password:", err)
	}

	user, err := store.CreateUser(context.Background(), CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	})
	if err != nil {
		log.Fatal("cannot create user:", err)
	}

	phrase_text := util.RandomString(12)
	phrase, err := store.CreatePhrase(context.Background(), CreatePhraseParams{
		Owner:     user.Username,
		State:     "draft",
		Phrase:    phrase_text,
		Author:    util.RandomOwner(),
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Fatal("cannot create phrase:", err)
	}

	require.Equal(t, phrase_text, phrase.Phrase)
	require.Equal(t, phrase.PublishedAt, sql.NullTime{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: false})
	require.Equal(t, phrase.State, "draft")
	require.NotEmpty(t, phrase.Phrase)
	issuedAt := time.Now()
	require.WithinDuration(t, phrase.CreatedAt, issuedAt, time.Second)
}

func TestUpdatePhrase(t *testing.T) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := NewStore(conn)

	hashedPassword, err := util.HashPassword(util.RandomString(12))
	if err != nil {
		log.Fatal("cannot create password:", err)
	}

	user, err := store.CreateUser(context.Background(), CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	})
	if err != nil {
		log.Fatal("cannot create user:", err)
	}

	phrase_text := util.RandomString(12)
	phrase, err := store.CreatePhrase(context.Background(), CreatePhraseParams{
		Owner:     user.Username,
		State:     "draft",
		Phrase:    phrase_text,
		Author:    util.RandomOwner(),
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Fatal("cannot create phrase:", err)
	}

	newAuthor := util.RandomOwner()
	newPhrase := util.RandomOwner()

	phraseUpdated, err := store.UpdatePhrase(context.Background(), UpdatePhraseParams{
		ID:     phrase.ID,
		Phrase: newPhrase,
		Author: newAuthor,
	})

	if err != nil {
		log.Fatal("cannot update phrase:", err)
	}

	require.Equal(t, phraseUpdated.Phrase, newPhrase)
	require.Equal(t, phraseUpdated.Author, newAuthor)
	require.Equal(t, phrase.PublishedAt, sql.NullTime{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: false})
	require.Equal(t, phrase.State, "draft")
	require.NotEmpty(t, phrase.Phrase)
	issuedAt := time.Now()
	require.WithinDuration(t, phrase.CreatedAt, issuedAt, time.Second)
}

func TestUpdatePhraseState(t *testing.T) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := NewStore(conn)

	hashedPassword, err := util.HashPassword(util.RandomString(12))
	if err != nil {
		log.Fatal("cannot create password:", err)
	}

	user, err := store.CreateUser(context.Background(), CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	})
	if err != nil {
		log.Fatal("cannot create user:", err)
	}

	phrase_text := util.RandomString(12)
	phrase, err := store.CreatePhrase(context.Background(), CreatePhraseParams{
		Owner:     user.Username,
		State:     "draft",
		Phrase:    phrase_text,
		Author:    util.RandomOwner(),
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Fatal("cannot create phrase:", err)
	}

	phraseUpdated, err := store.UpdatePhraseState(context.Background(), UpdatePhraseStateParams{
		ID:    phrase.ID,
		State: publishedPhraseState,
	})

	if err != nil {
		log.Fatal("cannot update phrase:", err)
	}

	require.Equal(t, phraseUpdated.State, publishedPhraseState)
	require.Equal(t, phrase.PublishedAt, sql.NullTime{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: false})
	require.Equal(t, phrase.State, "draft")
	require.NotEmpty(t, phrase.Phrase)
	issuedAt := time.Now()
	require.WithinDuration(t, phrase.CreatedAt, issuedAt, time.Second)
}

//func TestUpdatePublishedAt(t *testing.T) {
//
//	config, err := util.LoadConfig("../..")
//	if err != nil {
//		log.Fatal("cannot load config:", err)
//	}
//	conn, err := sql.Open(config.DBDriver, config.DBSource)
//	if err != nil {
//		log.Fatal("cannot connect to database:", err)
//	}
//	store := NewStore(conn)
//
//	hashedPassword, err := util.HashPassword(util.RandomString(12))
//	if err != nil {
//		log.Fatal("cannot create password:", err)
//	}
//
//	user, err := store.CreateUser(context.Background(), CreateUserParams{
//		Username:       util.RandomOwner(),
//		HashedPassword: hashedPassword,
//		FullName:       util.RandomOwner(),
//		Email:          util.RandomEmail(),
//	})
//	if err != nil {
//		log.Fatal("cannot create user:", err)
//	}
//
//	phrase_text := util.RandomString(12)
//	phrase, err := store.CreatePhrase(context.Background(), CreatePhraseParams{
//		Owner:     user.Username,
//		State:     "draft",
//		Phrase:    phrase_text,
//		Author:    util.RandomOwner(),
//		CreatedAt: time.Now(),
//	})
//
//	if err != nil {
//		log.Fatal("cannot create phrase:", err)
//	}
//
//	phraseUpdated, err := store.UpdatePhraseState(context.Background(), UpdatePhraseStateParams{
//		ID:    phrase.ID,
//		State: publishedPhraseState,
//	})
//
//	if err != nil {
//		log.Fatal("cannot update phrase:", err)
//	}
//
//	phraseUpdated, err = store.UpdatePublishedAt(context.Background(), phraseUpdated.ID)
//	if err != nil {
//		log.Fatal("cannot update published_at:", err)
//	}
//
//	publishedAt := time.Now()
//	require.WithinDuration(t, phraseUpdated.PublishedAt.Time, publishedAt, time.Second)
//}

func TestCountPublishedPhrasesToday(t *testing.T) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := NewStore(conn)

	hashedPassword, err := util.HashPassword(util.RandomString(12))
	if err != nil {
		log.Fatal("cannot create password:", err)
	}

	user, err := store.CreateUser(context.Background(), CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	})
	if err != nil {
		log.Fatal("cannot create user:", err)
	}

	phrase_text := util.RandomString(12)
	phrase, err := store.CreatePhrase(context.Background(), CreatePhraseParams{
		Owner:     user.Username,
		State:     "draft",
		Phrase:    phrase_text,
		Author:    util.RandomOwner(),
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Fatal("cannot create phrase:", err)
	}

	countPhrases, err := store.CountPhrasesPublishedToday(context.Background())
	if err != nil {
		log.Fatal("cannot count published prhases today:", err)
	}
	require.Equal(t, phrase_text, phrase.Phrase)
	require.Equal(t, countPhrases, int64(0))
}
