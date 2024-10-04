package services

import (
	"context"
	"fmt"
	db "github.com/cheojeg/top_phrases/db/sqlc"
	"strings"
)

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func formatMessage(phrase db.Phrase) string {
	phraseText := fmt.Sprintf("%s", escapeMarkdown(phrase.Phrase))
	if phrase.Author != "" {
		phraseText += fmt.Sprintf(" \\- *%s*", escapeMarkdown(phrase.Author))
	} else {
		phraseText += " \\- *Desconocido*"
	}
	return phraseText
}

func (s *Service) GetPhraseToPublish(ctx context.Context) (string, error) {
	phrase, _ := s.store.GetPhraseToPublish(ctx, 5)
	_, _ = s.store.UpdatePublishedAt(ctx, phrase.ID)
	fmt.Println(phrase)
	// TODO - Update last published at date
	//phraseToPublish := domain.Phrase{
	//	ID:        &phrase.ID,
	//	Phrase:    phrase.Phrase,
	//	Author:    phrase.Author,
	//	State:     phrase.State,
	//	CreatedAt: phrase.CreatedAt,
	//	// PublishedAt: phrase.PublishedAt,
	//}
	phraseText := formatMessage(phrase)

	return phraseText, nil
}
