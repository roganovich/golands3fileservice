CREATE TABLE "roles" (
     "id" serial PRIMARY KEY,
     "name" varchar NOT NULL
);

COMMENT ON COLUMN "roles"."name" IS 'Имя роли';

INSERT INTO roles (id, name) VALUES (1, 'Клиент');
INSERT INTO roles (id, name) VALUES (2, 'Капитан команды');
INSERT INTO roles (id, name) VALUES (10, 'Супер-администратор');
INSERT INTO roles (id, name) VALUES (11, 'Администратор');

ALTER TABLE users
    ADD COLUMN role_id int REFERENCES roles(id) DEFAULT 1;