CREATE TABLE "store_businesses" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "business_name" varchar NOT NULL
);
