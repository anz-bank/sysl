// Copyright 2016 The Sysl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package io.sysl;

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
