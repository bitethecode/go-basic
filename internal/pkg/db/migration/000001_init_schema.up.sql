CREATE TABLE "users" (
 "id" SERIAL PRIMARY KEY,
 "username" VARCHAR(255),
 "password" VARCHAR(255),
 "email" VARCHAR(255),
 "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
)
