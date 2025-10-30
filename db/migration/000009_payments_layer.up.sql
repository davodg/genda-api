CREATE TABLE "store_accounts" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "provider" varchar NOT NULL,
  "provider_account_id" varchar NOT NULL,
  "charges_enabled" boolean NOT NULL DEFAULT false,
  "payouts_enabled" boolean NOT NULL DEFAULT false,
  "default_currency" varchar DEFAULT 'BRL',
  "details" json,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  UNIQUE ("provider","provider_account_id"),
  CONSTRAINT fk_store_accounts_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE
);
CREATE INDEX ON "store_accounts" ("store_id");

CREATE TABLE "payments" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "appointment_id" uuid,
  "subscription_id" uuid,
  "provider" varchar,
  "provider_payment_id" varchar,
  "status" varchar NOT NULL DEFAULT 'processing' CHECK ("status" IN ('requires_payment_method','requires_confirmation','processing','succeeded','partially_refunded','refunded','failed','canceled')),
  "currency" varchar NOT NULL DEFAULT 'BRL',
  "amount_total" numeric(12,2) NOT NULL,
  "fee_platform_pct" numeric(5,2) NOT NULL DEFAULT 8.00,
  "fee_platform_amount" numeric(12,2),
  "fee_processing_amount" numeric(12,2),
  "amount_net" numeric(12,2),
  "captured_at" timestamp,
  "refunded_amount" numeric(12,2) DEFAULT 0,
  "metadata" json,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_payments_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE,
  CONSTRAINT fk_payments_user_id FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT,
  CONSTRAINT fk_payments_appointment_id FOREIGN KEY ("appointment_id") REFERENCES "store_appointments"("id") ON DELETE SET NULL,
  CONSTRAINT fk_payments_subscription_id FOREIGN KEY ("subscription_id") REFERENCES "subscriptions"("id") ON DELETE SET NULL
);
CREATE INDEX ON "payments" ("store_id");
CREATE INDEX ON "payments" ("user_id");

CREATE TABLE "refunds" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "payment_id" uuid NOT NULL,
  "provider_refund_id" varchar,
  "status" varchar NOT NULL DEFAULT 'pending' CHECK ("status" IN ('pending','succeeded','failed','canceled')),
  "amount" numeric(12,2) NOT NULL,
  "reason" varchar,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_refunds_payment_id FOREIGN KEY ("payment_id") REFERENCES "payments"("id") ON DELETE CASCADE
);
CREATE INDEX ON "refunds" ("payment_id");

CREATE TABLE "invoices" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "subscription_id" uuid,
  "provider_invoice_id" varchar,
  "period_start" timestamp,
  "period_end" timestamp,
  "currency" varchar NOT NULL DEFAULT 'BRL',
  "amount_total" numeric(12,2) NOT NULL,
  "fee_platform_amount" numeric(12,2) DEFAULT 0,
  "status" varchar NOT NULL DEFAULT 'open',
  "due_date" timestamp,
  "paid_at" timestamp,
  "payment_id" uuid,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  CONSTRAINT fk_invoices_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE,
  CONSTRAINT fk_invoices_subscription_id FOREIGN KEY ("subscription_id") REFERENCES "subscriptions"("id") ON DELETE SET NULL,
  CONSTRAINT fk_invoices_payment_id FOREIGN KEY ("payment_id") REFERENCES "payments"("id") ON DELETE SET NULL
);
CREATE INDEX ON "invoices" ("store_id");

CREATE TABLE "payouts" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "provider" varchar,
  "provider_transfer_id" varchar,
  "status" varchar NOT NULL DEFAULT 'pending' CHECK ("status" IN ('pending','in_transit','paid','failed','canceled')),
  "currency" varchar NOT NULL DEFAULT 'BRL',
  "amount" numeric(12,2) NOT NULL,
  "requested_at" timestamp NOT NULL DEFAULT now(),
  "paid_at" timestamp,
  "failure_code" varchar,
  "notes" varchar,
  CONSTRAINT fk_payouts_store_id FOREIGN KEY ("store_id") REFERENCES "stores"("id") ON DELETE CASCADE
);
CREATE INDEX ON "payouts" ("store_id");

CREATE TABLE "webhook_events" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "provider" varchar NOT NULL,
  "event_id" varchar NOT NULL,
  "event_type" varchar NOT NULL,
  "payload" json NOT NULL,
  "received_at" timestamp NOT NULL DEFAULT now(),
  "processed_at" timestamp,
  "status" varchar NOT NULL DEFAULT 'pending',
  UNIQUE ("provider","event_id")
);

-- add FK back from appointments to payments now that payments exists
ALTER TABLE "store_appointments"
  ADD CONSTRAINT fk_store_appointments_payment_id
  FOREIGN KEY ("payment_id") REFERENCES "payments"("id") ON DELETE SET NULL;
