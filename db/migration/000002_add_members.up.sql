CREATE TABLE "members" (
  "membername" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "traders" ADD FOREIGN KEY ("holder") REFERENCES "members" ("membername");

ALTER TABLE "traders" ADD CONSTRAINT "holder_currency_key" UNIQUE ("holder", "currency");