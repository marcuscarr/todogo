-- Database Definition
CREATE DATABASE marcus_demo;
\c marcus_demo

-- Table Definition
CREATE SEQUENCE IF NOT EXISTS todo_id_seq;
CREATE TABLE todos (
    "id" INTEGER NOT NULL DEFAULT nextval('todo_id_seq'::regclass),
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "status" BOOLEAN
    "created" TIMESTAMP(0),
    "modified" TIMESTAMP(0),
    PRIMARY KEY ("id")
);