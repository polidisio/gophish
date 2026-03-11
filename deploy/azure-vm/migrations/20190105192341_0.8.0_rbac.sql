-- +goose Up
CREATE TABLE "roles" (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255)
);

ALTER TABLE "users" ADD COLUMN "role_id" INTEGER;

CREATE TABLE "permissions" (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255)
);

CREATE TABLE "role_permissions" (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL
);

INSERT INTO "roles" ("slug", "name", "description")
VALUES
    ('admin', 'Admin', 'System administrator with full permissions'),
    ('user', 'User', 'User role with edit access to objects and campaigns');

INSERT INTO "permissions" ("slug", "name", "description")
VALUES
    ('view_objects', 'View Objects', 'View objects in Gophish'),
    ('modify_objects', 'Modify Objects', 'Create and edit objects in Gophish'),
    ('modify_system', 'Modify System', 'Manage system-wide configuration');

UPDATE "users" SET "role_id"=(
    SELECT "id" FROM "roles" WHERE "slug"='admin')
WHERE "id"=(
    SELECT "id" FROM "users" WHERE "username"='admin' OR "id"=(SELECT MIN("id") FROM "users") LIMIT 1);

UPDATE "users" SET "role_id"=(
    SELECT "id" FROM "roles" WHERE "slug"='user')
WHERE role_id IS NULL;

INSERT INTO "role_permissions" ("role_id", "permission_id")
SELECT r.id, p.id FROM roles AS r, "permissions" AS p
WHERE r.id IN (SELECT "id" FROM roles WHERE "slug"='admin' OR "slug"='user')
AND p.id=(SELECT "id" FROM "permissions" WHERE "slug"='view_objects');

INSERT INTO "role_permissions" ("role_id", "permission_id")
SELECT r.id, p.id FROM roles AS r, "permissions" AS p
WHERE r.id IN (SELECT "id" FROM roles WHERE "slug"='admin' OR "slug"='user')
AND p.id=(SELECT "id" FROM "permissions" WHERE "slug"='modify_objects');

INSERT INTO "role_permissions" ("role_id", "permission_id")
SELECT r.id, p.id FROM roles AS r, "permissions" AS p
WHERE r.id IN (SELECT "id" FROM roles WHERE "slug"='admin')
AND p.id=(SELECT "id" FROM "permissions" WHERE "slug"='modify_system');

-- +goose Down
DROP TABLE IF EXISTS "role_permissions";
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "roles";
ALTER TABLE "users" DROP COLUMN IF EXISTS "role_id";
