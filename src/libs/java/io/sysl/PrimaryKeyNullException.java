package io.sysl;

public class PrimaryKeyNullException extends RuntimeException {
    public PrimaryKeyNullException(String field) {
        super("PrimaryKeyNullException(\"" + field + "\")");
    }
}