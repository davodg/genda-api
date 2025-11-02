CREATE TABLE "stores" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "name" varchar NOT NULL,
  "owner_id" uuid NOT NULL,
  "type" varchar NOT NULL,
  "location" json,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_stores_owner_id FOREIGN KEY ("owner_id") REFERENCES "users"("id") ON DELETE RESTRICT,
);
CREATE INDEX ON "stores" ("owner_id");
