-- name: CreatePhrase :one
INSERT INTO phrases (
    owner, state, phrase, author, created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdatePhraseState :one
UPDATE phrases
SET state = $2
WHERE id = $1
RETURNING *;

-- name: UpdatePhrase :one
UPDATE phrases
SET phrase = $2, author = $3
WHERE id = $1
RETURNING *;

-- name: GetPhraseToPublish :one
SELECT *
FROM phrases
WHERE published_at IS NULL OR published_at < NOW() - INTERVAL '15 days'
ORDER BY RANDOM()
LIMIT 1;

-- name: UpdatePublishedAt :many
UPDATE phrases
SET published_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListPhrases :many
SELECT * FROM phrases;

-- name: GetPhraseByID :one
SELECT *
FROM phrases
WHERE id = $1;

-- name: CountDraftPhrases :one
SELECT COUNT(*)
FROM phrases
WHERE state = 'draft';