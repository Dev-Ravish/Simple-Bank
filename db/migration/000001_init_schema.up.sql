CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "transfer_from_account" bigint NOT NULL,
  "transfer_to_account" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "transactions" ("account_id");

CREATE INDEX ON "transfers" ("transfer_from_account");

CREATE INDEX ON "transfers" ("transfer_to_account");

CREATE INDEX ON "transfers" ("transfer_from_account", "transfer_to_account");

COMMENT ON COLUMN "transactions"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "transactions" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("transfer_from_account") REFERENCES "account" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("transfer_to_account") REFERENCES "account" ("id");