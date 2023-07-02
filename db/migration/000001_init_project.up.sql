CREATE TABLE "traders" (
  "id" bigserial PRIMARY KEY,
  "account" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "records" (
  "id" bigserial PRIMARY KEY,
  "trader_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "from_trader_id" bigint NOT NULL,
  "to_trader_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "traders" ("account");

CREATE INDEX ON "records" ("trader_id");

CREATE INDEX ON "payments" ("from_trader_id");

CREATE INDEX ON "payments" ("to_trader_id");

CREATE INDEX ON "payments" ("from_trader_id", "to_trader_id");

COMMENT ON COLUMN "records"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "payments"."amount" IS 'most be positive';

ALTER TABLE "records" ADD FOREIGN KEY ("trader_id") REFERENCES "traders" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("from_trader_id") REFERENCES "traders" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("to_trader_id") REFERENCES "traders" ("id");
