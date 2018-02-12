package io.sysl;

import java.lang.Comparable;

import java.math.BigDecimal;

public class ExprHelper {

    public static boolean areEqual(Object a, Object b) {
        return a != null ? a.equals(b) : b == null;
    }

    public static boolean areEqual(BigDecimal a, BigDecimal b) {
        return doCompare(a, b) == 0;
    }

    public static <T> int doCompare(Comparable<T> a, T b) {
        return a != null && b != null ? a.compareTo(b) :
            a == null && b == null ? 0 :
            a == null ? 1 : -1;
    }

    public static Integer toInteger(String s) {
        if (s == null) {
            return null;
        }
        return Integer.parseInt(s);
    }

    public static Integer toInteger(String s, Integer fallback) {
        try {
            return toInteger(s);
        } catch (Exception ex) {
            return fallback;
        }
    }

    public static Integer minus(Integer x) {
        return -x;
    }

    public static Double minus(Double x) {
        return -x;
    }

    public static BigDecimal minus(BigDecimal x) {
        return x.negate();
    }

    public static Integer plus(Integer x) {
        return +x;
    }

    public static Double plus(Double x) {
        return +x;
    }

    public static BigDecimal plus(BigDecimal x) {
        return x.plus();
    }

}
