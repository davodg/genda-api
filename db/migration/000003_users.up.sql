CREATE TABLE "users" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "is_admin" boolean NOT NULL DEFAULT false,
  "type" varchar NOT NULL,
  "photo_url" varchar,
  "birthdate" date NOT NULL,
  "social_number" varchar,
  "gender" varchar,
  "phone_number" varchar,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now()
);
