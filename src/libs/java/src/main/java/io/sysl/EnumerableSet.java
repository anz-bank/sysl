package io.sysl;

import java.util.HashSet;

public class EnumerableSet<T> extends Enumerable<T> {
    public EnumerableSet(Iterable<T> input) {
        for (T e : input) {
            data.add(e);
        }
    }

    public boolean contains(T t) {
        return data.contains(t);
    }

    public Enumerator<T> enumerator() {
        return Enumeration.enumerator(data.iterator());
    }

    private HashSet<T> data = new HashSet<T>();
}
