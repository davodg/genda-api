CREATE TABLE "store_plans" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "price" numeric(12,2) NOT NULL,
  "currency" varchar NOT NULL DEFAULT 'BRL',
  "plan_type" varchar,
  "frequency" varchar, -- 'day','week','month','year'
  "metadata" jsonb,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_store_plans_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE
);
CREATE INDEX ON "store_plans" ("store_id");

CREATE TABLE "subscriptions" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "store_id" uuid NOT NULL,
  "store_plan_id" uuid NOT NULL,
  "status" varchar NOT NULL DEFAULT 'active' CHECK ("status" IN ('trialing','active','past_due','paused','canceled','unpaid')),
  "current_period_start" timestamp,
  "current_period_end" timestamp,
  "cancel_at" timestamp,
  "canceled_at" timestamp,
  "trial_end" timestamp,
  "provider" varchar,
  "provider_sub_id" varchar,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_subscriptions_user_id FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
  CONSTRAINT fk_subscriptions_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE,
  CONSTRAINT fk_subscriptions_store_plan_id FOREIGN KEY ("store_plan_id") REFERENCES "store_plans"("id") ON DELETE RESTRICT
);
CREATE INDEX ON "subscriptions" ("user_id");
CREATE INDEX ON "subscriptions" ("store_id");
CREATE INDEX ON "subscriptions" ("store_plan_id");
