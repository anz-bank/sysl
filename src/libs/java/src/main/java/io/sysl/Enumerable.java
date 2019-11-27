package io.sysl;

import java.lang.Comparable;

import java.math.BigDecimal;

import java.util.ArrayList;

public abstract class Enumerable<T> implements java.lang.Iterable<T> {
    public abstract Enumerator<T> enumerator();

    public java.util.Iterator<T> iterator() {
        final Enumerator<T> enumerator = this.enumerator();

        return new java.util.Iterator<T>() {
            @Override
            public boolean hasNext() {
                return have;
            }

            @Override
            public T next() throws java.util.NoSuchElementException {
                if (!have) {
                    throw new java.util.NoSuchElementException();
                }
                T result = enumerator.current();
                have = enumerator.moveNext();
                return result;
            }

            @Override
            public void remove()
                    throws java.lang.UnsupportedOperationException {
                throw new java.lang.UnsupportedOperationException();
            }

            private boolean have = enumerator.moveNext();
        };
    }

    public boolean isEmpty() {
        return !iterator().hasNext();
    }

    public int size() {
        return sizeWithLimit(Integer.MAX_VALUE);
    }

    public Integer count() {
        return size();
    }

    public int sizeWithLimit(int limit) {
        int n = 0;
        for (T item : this) {
            if (++n == limit) {
                break;
            }
        }
        return n;
    }

    public ArrayList<T> asArrayList() {
        ArrayList<T> result = new ArrayList<T>();
        for (T e : this) {
            result.add(e);
        }
        return result;
    }

    public T singleOrNull() throws RuntimeException {
        T result = null;
        for (T item : this) {
            if (result != null) {
                throw new SyslException("size() == " + size() + " > 1");
            }
            result = item;
        }
        return result;
    }

    public T single() throws RuntimeException {
        T result = singleOrNull();
        if (result == null) {
            throw new SyslException("size() == " + size() + " != 1");
        }
        return result;
    }

    public <U> Enumerable<U> map(Expr<U, T> map) {
        return Enumeration.map(this, map);
    }

    public Enumerable<T> orderBy(java.util.Comparator<T> comp) {
        return Enumeration.orderBy(this, comp);
    }

    public interface Reducer<T, U> {
        public U reduce(U u, T t);
    }

    public <U> U reduce(U u, Reducer<T, U> reducer) {
        for (T t : this) {
            u = reducer.reduce(u, t);
        }
        return u;
    }

    // public <U> U reduce(Reducer<T, U> reducer) {
    //     Enumerator<T> source = enumerator();
    //     if (!source.moveNext()) {
    //         throw new java.lang.RuntimeException("Nothing to reduce");
    //     }
    //     U u = source.current();
    //     while (source.moveNext()) {
    //         u = reducer.reduce(u, source.current());
    //     }
    //     return u;
    // }

    public <U extends Comparable<U>> U min(final Expr<U, T> expr) {
        for (T e : this) {
            return reduce(expr.evaluate(e), new Reducer<T, U>() {
                @Override
                public U reduce(U i, T t) {
                    U j = expr.evaluate(t);
                    return i.compareTo(j) < 0 ? i : j;
                }
            });
        }
        return null;
    }

    public <U extends Comparable<U>> U max(final Expr<U, T> expr) {
        for (T e : this) {
            return reduce(expr.evaluate(e), new Reducer<T, U>() {
                @Override
                public U reduce(U i, T t) {
                    U j = expr.evaluate(t);
                    return i.compareTo(j) < 0 ? j : i;
                }
            });
        }
        return null;
    }

    public Integer sumInteger(final Expr<Integer, T> expr) {
        return reduce(0, new Reducer<T, Integer>() {
            @Override
            public Integer reduce(Integer i, T t) {
                return i + expr.evaluate(t);
            }
        });
    }

    public Double sumDouble(final Expr<Double, T> expr) {
        return reduce(0.0, new Reducer<T, Double>() {
            @Override
            public Double reduce(Double i, T t) {
                return i + expr.evaluate(t);
            }
        });
    }

    public BigDecimal sumBigDecimal(final Expr<BigDecimal, T> expr) {
        return reduce(BigDecimal.ZERO, new Reducer<T, BigDecimal>() {
            @Override
            public BigDecimal reduce(BigDecimal i, T t) {
                return i.add(expr.evaluate(t));
            }
        });
    }

    // private static class Pair<T> {
    //     T a;
    //     T b;
    //     Pair(T a, T b) {
    //         this.a = a;
    //         this.b = b;
    //     }
    // }

    // public average(Expr<Double, T> expr) {
    //     Pair<Double> p = reduce(
    //         new Pair<Double>(0, 0), new Reducer<T, Pair<Double>>() {
    //             @Override
    //             public Pair<Double> reduce(Pair<Double> i, T t) {
    //                 return new Pair<Double>(i.a + expr.evaluate(t), i.b + 1);
    //             }
    //         });
    //     return p.a / p.b;
    // }

    // public average(Expr<BigDecimal, T> expr) {
    //     Pair<BigDecimal> p = reduce(
    //         new Pair<BigDecimal>(0, 0), new Reducer<T, Pair<BigDecimal>>() {
    //             @Override
    //             public Pair<BigDecimal> reduce(Pair<BigDecimal> i, T t) {
    //                 return new Pair<BigDecimal>(
    //                     i.a.add(expr.evaluate(t)), i.b + 1);
    //             }
    //         });
    //     return p.a / p.b;
    // }

}
