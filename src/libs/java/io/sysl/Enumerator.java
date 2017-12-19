package io.sysl;

public interface Enumerator<T> {
    public boolean moveNext();
    public T current();
}
