ALTER TABLE "store_appointments" DROP CONSTRAINT IF EXISTS fk_store_appointments_payment_id;

DROP TABLE IF EXISTS "webhook_events";
DROP TABLE IF EXISTS "payouts";
DROP TABLE IF EXISTS "invoices";
DROP TABLE IF EXISTS "refunds";
DROP TABLE IF EXISTS "payments";
DROP TABLE IF EXISTS "store_accounts";
