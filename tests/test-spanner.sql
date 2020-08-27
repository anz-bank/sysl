CREATE DATABASE customeraccounts;

CREATE TABLE Account (
    AccountNum STRING(23) NOT NULL,
    BSB        STRING(6)  NOT NULL,
    Balance    INT64      NOT NULL,
) PRIMARY KEY (AccountNum);

CREATE INDEX AccountsByNum ON Account (AccountNum DESC);

CREATE TABLE Customer (
    CustomerID STRING(36)  NOT NULL,
    FirstName  STRING(64)  NOT NULL,
    LastName   STRING(64)  NOT NULL,
    Email      STRING(256),
    Mobile     STRING(10),
    NetWorth   FLOAT64,
) PRIMARY KEY (CustomerID);

CREATE UNIQUE NULL_FILTERED INDEX CustomerByEmail ON Customer (Email ASC, Mobile DESC);

CREATE TABLE CustomerHasAccount (
    CustomerID  STRING(36) NOT NULL,
    AccountNum  STRING(23) NOT NULL,
    LegalRole   STRING(10) NOT NULL,
    BranchID    STRING(6)  NOT NULL,
    Permissions ARRAY<STRING(10)>,
    CONSTRAINT FK_CustomerID FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
    CONSTRAINT FK_AccountNum FOREIGN KEY (AccountNum, BranchID) REFERENCES Account (AccountNum, BSB),
) PRIMARY KEY (AccountNum ASC, CustomerID);

CREATE TABLE AccountAddress (
    AccountNum        STRING(23) NOT NULL,
    AddressPostCode   STRING(10) NOT NULL,
    LastUpdated       TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    AddressLine1      BYTES(MAX),
    AddressLine2      STRING(0x100),
    AddressLine3      STRING(100),
) PRIMARY KEY (AccountNum ASC, AddressPostCode DESC),
INTERLEAVE IN PARENT Account ON DELETE CASCADE;
