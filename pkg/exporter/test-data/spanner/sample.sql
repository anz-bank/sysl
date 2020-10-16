CREATE TABLE Types (
    array ARRAY<INT64> NOT NULL,
    array_set ARRAY<INT64> NOT NULL,
    bytes_ BYTES NOT NULL,
    bytes_1 BYTES NOT NULL,
    bytes_1024_hex BYTES NOT NULL,
    bytes_max BYTES NOT NULL,
    bool_ BOOL NOT NULL,
    opt BOOL,
    int64_ INT64 NOT NULL,
    float64_ FLOAT64 NOT NULL,
    date_ DATE NOT NULL,
    timestamp_ TIMESTAMP NOT NULL,
    timestamp_commit TIMESTAMP NOT NULL,
    string_ STRING NOT NULL,
    string_1 STRING(1) NOT NULL,
    string_100 STRING(100) NOT NULL,
    string_max STRING NOT NULL,
    cust_id INT64 NOT NULL,
    FOREIGN KEY (cust_id) REFERENCES Customer (id)
) PRIMARY KEY (array);
CREATE INDEX Ix ON Types(float64_, int64_ DESC);

CREATE TABLE Customer (
    id INT64 NOT NULL,
    ref_desc INT64 NOT NULL,
    ref_asc INT64 NOT NULL,
) PRIMARY KEY (id, ref_desc, ref_asc);

CREATE TABLE Constraints (
    cust_id INT64 NOT NULL,
    types_int INT64 NOT NULL,
    types_string STRING NOT NULL,
    FOREIGN KEY (cust_id) REFERENCES Customer (id),
    FOREIGN KEY (types_int) REFERENCES Types (int64_),
    FOREIGN KEY (types_string) REFERENCES Types (string_)
);
