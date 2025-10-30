CREATE TABLE "store_appointments" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "start_at" timestamp NOT NULL,
  "end_at" timestamp NOT NULL,
  "status" varchar NOT NULL DEFAULT 'pending' CHECK ("status" IN ('pending','confirmed','completed','canceled','no_show')),
  "hold_expires_at" timestamp,
  "price" numeric(12,2),
  "currency" varchar DEFAULT 'BRL',
  "fee_platform" numeric(12,2),
  "payment_id" uuid,
  "notes" varchar,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_store_appointments_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE,
  CONSTRAINT fk_store_appointments_user_id FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT
);
CREATE INDEX ON "store_appointments" ("store_id");
CREATE INDEX ON "store_appointments" ("user_id");
CREATE INDEX ON "store_appointments" ("start_at","end_at");

ALTER TABLE "store_appointments"
  ADD CONSTRAINT no_overlap_per_store
  EXCLUDE USING gist (
    "store_id" WITH =,
    tsrange("start_at","end_at",'[)') WITH &&
  )
  WHERE ("status" IN ('pending','confirmed'));
