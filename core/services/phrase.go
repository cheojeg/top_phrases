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

func (s *Service) GetPhraseToPublish(ctx context.Context, days int64) (string, error) {
	phrase, err := s.store.GetPhraseToPublish(ctx, days)
	if err != nil {
		return "", err
	}
	_, err = s.store.UpdatePublishedAt(ctx, phrase.ID)
	if err != nil {
		return "", err
	}
	fmt.Println(phrase)
	phraseText := formatMessage(phrase)
	return phraseText, nil
}
