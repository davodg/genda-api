CREATE TABLE "store_ratings" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "rating" int NOT NULL CHECK ("rating" BETWEEN 1 AND 5),
  "message" varchar,
  "created_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_store_ratings_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE,
  CONSTRAINT fk_store_ratings_user_id FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);
CREATE INDEX ON "store_ratings" ("store_id");
CREATE INDEX ON "store_ratings" ("user_id");

CREATE UNIQUE INDEX "uniq_daily_review" ON "store_ratings"
("store_id","user_id",(date_trunc('day',"created_at")));
