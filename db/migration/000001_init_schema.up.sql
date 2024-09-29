CREATE TABLE "users" (
     "username" varchar PRIMARY KEY,
     "hashed_password" varchar NOT NULL,
     "full_name" varchar NOT NULL,
     "email" varchar UNIQUE NOT NULL,
     "password_changer_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
     "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "phrases" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "state" varchar NOT NULL,
    "phrase" varchar NOT NULL,
    "author" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "published_at" timestamptz
);

CREATE TABLE "phrase_states" (
 "state" varchar PRIMARY KEY
);

CREATE INDEX ON "phrases" ("owner");

CREATE INDEX ON "phrases" ("state");

ALTER TABLE "phrases" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "phrases" ADD FOREIGN KEY ("state") REFERENCES "phrase_states" ("state");