-- Database Definition
CREATE DATABASE marcus_demo;
\c marcus_demo

-- Table Definition
CREATE SEQUENCE IF NOT EXISTS todo_id_seq;
CREATE TABLE todos (
    "id" INTEGER NOT NULL DEFAULT nextval('todo_id_seq'::regclass),
    "title" TEXT,
    "description" TEXT,
    "status" BOOLEAN DEFAULT FALSE,
    "created" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP(0),
    "modified" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP(0),
    PRIMARY KEY ("id")
);