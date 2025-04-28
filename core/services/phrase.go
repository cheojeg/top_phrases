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
		"—", "\\—",
	)
	return replacer.Replace(text)
}

func formatMessageMarkdown(phrase db.Phrase) string {
	phraseText := fmt.Sprintf("%s", escapeMarkdown(phrase.Phrase))
	if phrase.Author != "" {
		phraseText += fmt.Sprintf(" \\— *%s*", escapeMarkdown(phrase.Author))
	} else {
		phraseText += " \\— *Desconocido*"
	}
	return phraseText
}

func formatMessage(phrase db.Phrase) string {
	phraseText := fmt.Sprintf("%s", phrase.Phrase)
	if phrase.Author != "" {
		phraseText += fmt.Sprintf(" — %s", phrase.Author)
	} else {
		phraseText += " — Desconocido"
	}
	return phraseText
}

func (s *Service) GetPhraseToPublish(ctx context.Context, days int64) (string, error) {
	countPhrases, err := s.store.CountPhrasesPublishedToday(context.Background())
	if err != nil {
		return "", err
	}

	if countPhrases > 0 {
		return "", fmt.Errorf("A quote was already published today")
	}

	phrase, err := s.store.GetPhraseToPublish(ctx, days)
	if err != nil {
		return "", err
	}
	_, err = s.store.UpdatePublishedAt(ctx, phrase.ID)
	if err != nil {
		return "", err
	}
	phraseText := formatMessage(phrase)
	return phraseText, nil
}
