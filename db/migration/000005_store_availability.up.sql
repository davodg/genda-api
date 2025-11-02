CREATE TABLE "store_availability" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "availability" json NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_store_availability_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE
);
CREATE INDEX ON "store_availability" ("store_id");
