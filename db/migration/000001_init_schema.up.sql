CREATE EXTENSION pgcrypto;
CREATE TABLE "users" (
                         "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
                         "user_name" varchar NOT NULL,
                         "email" varchar NOT NULL,
                         "password" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now()),
                         "updated_at" timestamptz NOT NULL DEFAULT (now()),
                         "verified_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_role" (
                             "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
                             "user_id" uuid NOT NULL,
                             "role_id" int NOT NULL,
                             "created_at" timestamptz NOT NULL DEFAULT (now()),
                             "updated_at" timestamptz NOT NULL DEFAULT (now()),
                             "verified_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_profile" (
                                "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
                                "user_id" uuid NOT NULL,
                                "first_name" varchar,
                                "last_name" varchar,
                                "business_name" varchar,
                                "street_address" varchar,
                                "city" varchar,
                                "state" varchar,
                                "zip" varchar,
                                "country_code" int,
                                "created_at" timestamptz NOT NULL DEFAULT (now()),
                                "updated_at" timestamptz NOT NULL DEFAULT (now()),
                                "verified_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "users" ("user_name");

CREATE INDEX ON "user_role" ("user_id");

CREATE INDEX ON "user_role" ("role_id");

CREATE INDEX ON "user_profile" ("user_id");

ALTER TABLE "user_role" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_profile" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");