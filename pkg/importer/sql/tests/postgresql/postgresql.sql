CREATE DATABASE customeraccounts;

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('foo_bar', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';

CREATE TABLE PayID (
    PayID   VARCHAR(23) NOT NULL PRIMARY KEY,
    BSB     VARCHAR(6)  NOT NULL,
    float   FLOAT       NOT NULL,
);

CREATE TABLE Account (
    AccountNum   VARCHAR(23) NOT NULL,
    BSB          VARCHAR(6)  NOT NULL,
    Balance      BIGINT      NOT NULL,
    CreationDate DATE        NOT NULL,
    Table        Varchar(32),
    PRIMARY KEY (AccountNum)
);

CREATE INDEX AccountsByNum ON Account (AccountNum DESC);
CREATE INDEX Complex ON Account (AccountNum, BSB DESC, Balance ASC);

CREATE TABLE Customer (
    CustomerID VARCHAR(36)  NOT NULL,
    FirstName  VARCHAR(64)  NOT NULL,
    LastName   VARCHAR(64)  NOT NULL,
    Email      VARCHAR(256),
    Mobile     VARCHAR(10),
    NetWorth   NUMERIC,
    Int        INT          NOT NULL,
    PRIMARY KEY (CustomerID)
);

CREATE UNIQUE INDEX CustomerByEmail ON Customer (Email, Mobile DESC);

CREATE TABLE CustomerHasAccount (
    CustomerID  VARCHAR(36) NOT NULL,
    Customer    VARCHAR(36) NOT NULL,
    AccountNum  VARCHAR(23) NOT NULL,
    LegalRole   VARCHAR(10) NOT NULL,
    BranchID    VARCHAR(6)  NOT NULL,
    Permissions VARCHAR(10)[],
    PRIMARY KEY (AccountNum, CustomerID),
    FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
    CONSTRAINT FK_AccountNum FOREIGN KEY (AccountNum, BranchID) REFERENCES Account (AccountNum, BSB),
);

CREATE TABLE AccountAddress (
    AccountNum      VARCHAR(23) NOT NULL PRIMARY KEY,
    AddressPostCode VARCHAR(10) NOT NULL PRIMARY KEY,
    LastUpdated     TIMESTAMP,
    AddressLine1    BYTES(MAX),
    AddressLine2    VARCHAR(0x100),
    AddressLine3    BYTES(100),
    PRIMARY KEY (AccountNum, AddressPostCode)
);

CREATE SEQUENCE AccountNumSeq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE Account OWNER TO local;
ALTER SEQUENCE AccountNumSeq OWNED BY Account.AccountNum;
ALTER TABLE ONLY Account ALTER COLUMN AccountNum SET DEFAULT nextval('AccountNumSeq'::regclass);
