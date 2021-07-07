CREATE DATABASE test;

CREATE TABLE IF NOT EXISTS namespace.t (
    foo ARRAY<STRUCT<
        f STRING,
        bar ARRAY<STRUCT<
            b STRING NOT NULL
        >>
    >>
)
