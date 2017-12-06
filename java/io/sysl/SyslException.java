package io.sysl;

import java.lang.RuntimeException;

public class SyslException extends RuntimeException {

    public SyslException() { }

    public SyslException(String message) {
        super(message);
    }

}
