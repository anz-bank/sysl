package io.sysl;

import java.lang.Iterable;
import java.lang.UnsupportedOperationException;

import java.util.Collections;
import java.util.Comparator;
import java.util.Iterator;
import java.util.HashSet;
import java.util.NoSuchElementException;

public class Enumeration {

    public static <T> Enumerable<T> enumerable(final Iterable<T> iterable) {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                return Enumeration.enumerator(iterator());
            }

            @Override
            public Iterator<T> iterator() {
                return iterable.iterator();
            }
        };
    }

    public static <T> Enumerator<T> enumerator(final Iterator<T> iterator) {
        return new Enumerator<T>() {
            @Override
            public boolean moveNext() {
                if (iterator.hasNext()) {
                    value = iterator.next();
                    return true;
                }
                return false;
            }

            @Override
            public T current() {
                return value;
            }

            T value;
        };
    }

    public static <T> Enumerable<T> dedup(final Enumerable<T> input) {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                final Enumerator<T> source = input.enumerator();

                return new Enumerator<T>() {
                    @Override
                    public boolean moveNext() {
                        while (source.moveNext()) {
                            if (!seen.contains(current())) {
                                seen.add(current());
                                return true;
                            }
                        }
                        return false;
                    }

                    @Override
                    public T current() {
                        return source.current();
                    }

                    private HashSet<T> seen = new HashSet<T>();
                };
            }
        };
    }

    public static <T> Enumerable<T> any(
            final Enumerable<T> input, final int n) {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                final Enumerator<T> source = input.enumerator();

                return new Enumerator<T>() {
                    @Override
                    public boolean moveNext() {
                        if (i != 0 && source.moveNext()) {
                            i--;
                            return true;
                        }
                        return false;
                    }

                    @Override
                    public T current() {
                        return source.current();
                    }

                    private int i = n;
                };
            }
        };
    }

    public static <T> Enumerable<T> empty() {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                return new Enumerator<T>() {
                    @Override
                    public boolean moveNext() {
                        return false;
                    }

                    @Override
                    public T current() {
                        return null;
                    }
                };
            }
        };
    }

    public static <T> Enumerable<T> where(
            final Enumerable<T> input, final Expr<Boolean, T> pred) {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                final Enumerator<T> source = input.enumerator();

                return new Enumerator<T>() {
                    @Override
                    public boolean moveNext() {
                        while (source.moveNext()) {
                            if (pred.evaluate(current())) {
                                return true;
                            }
                        }
                        return false;
                    }

                    @Override
                    public T current() {
                        return source.current();
                    }
                };
            }
        };
    }

    public static <T, U> Enumerable<U> map(
            final Enumerable<T> input, final Expr<U, T> map) {
        return new Enumerable<U>() {
            @Override
            public Enumerator<U> enumerator() {
                final Enumerator<T> source = input.enumerator();

                return new Enumerator<U>() {
                    @Override
                    public boolean moveNext() {
                        if (source.moveNext()) {
                            curr = map.evaluate(source.current());
                            return true;
                        }
                        curr = null;
                        return false;
                    }

                    @Override
                    public U current() {
                        return curr;
                    }

                    private U curr;
                };
            }
        };
    }

    public static <T, U> Enumerable<U> flatten(
            final Enumerable<T> input, final Expr<Enumerable<U>, T> nested) {
        return dedup(new Enumerable<U>() {
            @Override
            public Enumerator<U> enumerator() {
                final Enumerator<T> outer = input.enumerator();

                return new Enumerator<U>() {
                    @Override
                    public boolean moveNext() {
                        while (inner != null && !inner.moveNext()) {
                            if ((inner = nextOuter()) == null) {
                                return false;
                            }
                        }
                        return inner != null;
                    }

                    @Override
                    public U current() {
                        return inner.current();
                    }

                    private Enumerator<U> nextOuter() {
                        if (outer.moveNext()) {
                            return nested.evaluate(outer.current()).enumerator();
                        }
                        return null;
                    }

                    private Enumerator<U> inner = nextOuter();
                };
            }
        });
    }

    public static <T> Enumerable<T> orderBy(
            final Enumerable<T> input, final Comparator<T> comp) {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                // Default initial capacity from https://goo.gl/RORGLq
                final java.util.PriorityQueue<T> pq =
                    new java.util.PriorityQueue<T>(11, comp);

                for (T item : input) {
                    pq.add(item);
                }

                return new Enumerator<T>() {
                    @Override
                    public boolean moveNext() {
                        item = pq.poll();
                        return item != null;
                    }

                    @Override
                    public T current() {
                        return item;
                    }

                    private T item;
                };
            }
        };
    }

    public static <T> Enumerable<T> first(
            final Enumerable<T> input, final int n,
            final Comparator<T> comp) {
        return new Enumerable<T>() {
            @Override
            public Enumerator<T> enumerator() {
                // Default from https://goo.gl/RORGLq
                final java.util.PriorityQueue<T> pq =
                    new java.util.PriorityQueue<T>(
                        n + 1, Collections.reverseOrder(comp));

                for (T item : input) {
                    pq.add(item);
                    if (pq.size() == n + 1) {
                        pq.poll();
                    }
                }

                return Enumeration.enumerator(pq.iterator());
            }
        };
    }

    public interface Ranker<T, R> extends Comparator<T> {
        public R ranked(T t, int r);
    }

    public static <T, U> Enumerable<U> rank(
            final Enumerable<T> input, final Ranker<T, U> ranker) {
        return new Enumerable<U>() {
            @Override
            public Enumerator<U> enumerator() {
                final Enumerator<T> source =
                    Enumeration.orderBy(input, ranker).enumerator();

                return new Enumerator<U>() {
                    @Override
                    public boolean moveNext() {
                        if (source.moveNext()) {
                            T curr_in = source.current();
                            int r = (
                                prev_in == null ? -1 :
                                ranker.compare(prev_in, curr_in));
                            assert r <= 0: curr_in;
                            if (r != 0) {
                                rank = count;
                                prev_in = curr_in;
                            }
                            ++count;
                            curr = ranker.ranked(source.current(), rank);
                            return true;
                        }
                        curr = null;
                        return false;
                    }

                    @Override
                    public U current() {
                        return curr;
                    }

                    private T prev_in;
                    private U curr;
                    private int rank = 0;
                    private int count = 0;
                };
            }
        };
    }

}
