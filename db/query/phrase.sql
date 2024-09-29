-- name: CreatePhrase :one
INSERT INTO phrases (
    owner, state, phrase, author, created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdatePhraseState :one
UPDATE phrases
SET state = $2, published_at = $3
WHERE id = $1
RETURNING *;