CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "billingdata"(
	"org_id"  varchar NOT NULL,
	"activity_id"  varchar NOT NULL,
	"candidate_id"  varchar NOT NULL,
	"action"  varchar NOT NULL,
	"date"  timestamp NOT NULL (now())
);
CREATE TABLE "customerdata"(
	"org_id"  varchar NOT NULL,
	"customer_id"  varchar NOT NULL,
	"name"  varchar NOT NULL,
	"email"  varchar NOT NULL,
	"created_at"  timestamp NOT NULL (now())
);