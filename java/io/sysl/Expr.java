package io.sysl;

public interface Expr<T, E> {

    public abstract T evaluate(E entity);

}
