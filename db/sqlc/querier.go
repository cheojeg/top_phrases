// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreatePhrase(ctx context.Context, arg CreatePhraseParams) (Phrase, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetPhraseToPublish(ctx context.Context) (Phrase, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, username string) (User, error)
	UpdatePhrase(ctx context.Context, arg UpdatePhraseParams) (Phrase, error)
	UpdatePhraseState(ctx context.Context, arg UpdatePhraseStateParams) (Phrase, error)
}

var _ Querier = (*Queries)(nil)
