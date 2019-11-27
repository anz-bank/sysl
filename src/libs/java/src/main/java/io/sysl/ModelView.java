package io.sysl;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import java.lang.System;

import java.nio.charset.StandardCharsets;

import java.util.HashMap;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.joda.time.DateTime;
import org.joda.time.LocalDate;

public class ModelView {
    public static class EmptyTuple {
        public static EmptyTuple theOne = new EmptyTuple();

        public static class View extends Enumerable<EmptyTuple> {
            @Override
            public Enumerator<EmptyTuple> enumerator() {
                return Enumeration.<EmptyTuple>empty().enumerator();
            }
        }

        private EmptyTuple() { }
    }

    protected int autoinc() {
        return autoinc("");
    }

    protected int autoinc(String key) {
        Integer i = autoincs.get(key);
        if (i == null) {
            i = 1;
        } else {
            i += 1;
        }
        autoincs.put(key, i);
        return i;
    }

    protected String concat(String[] params) {
        return concat(params, "");
    }

    protected String concat(String[] params, String sep) {
        StringBuilder sb = new StringBuilder();
        boolean added = false;
        for (String param : params) {
            if (param != null) {
                if (added) {
                    sb.append(sep);
                } else {
                    added = true;
                }
                sb.append(param);
            }
        }
        return added ? sb.toString() : null;
    }

    private int firstNotOf(String s, String toks) {
        int len = s.length();
        int i = 0;
        for (; i < len; ++i) {
            if (toks.indexOf(s.charAt(i)) != -1) {
                break;
            }
        }
        return i;
    }

    private int lastNotOf(String s, String toks) {
        int i = s.length();
        while (i-- != 0) {
            if (toks.indexOf(s.charAt(i)) != -1) {
                break;
            }
        }
        return i;
    }

    protected Boolean log(String fmt, Object... params) {
        System.out.format(fmt, params);
        return true;
    }

    protected String lstrip(String s, String toks) {
        return s == null ? s : s.substring(firstNotOf(s, toks));
    }

    protected String rstrip(String s, String toks) {
        return s == null ? s : s.substring(0, lastNotOf(s, toks) + 1);
    }

    protected String strip(String s, String toks) {
        if (s == null) {
            return null;
        }
        int start = firstNotOf(s, toks);
        int finish = lastNotOf(s, toks) + 1;
        if (finish <= start) {
            return "";
        }
        return s.substring(start, finish);
    }

    protected String regsub(String regex, String replacement, String text) {
        Pattern p = patterns.get(regex);
        if (p == null) {
            p = Pattern.compile(regex);
            patterns.put(regex, p);
        }
        Matcher m = p.matcher(text);
        return m.replaceAll(replacement);
    }

    protected int clamp(int value, int min, int max) {
        return value < min ? min : value > max ? max : value;
    }

    protected double clamp(double value, double min, double max) {
        return value < min ? min : value > max ? max : value;
    }

    protected java.math.BigDecimal clamp(
            java.math.BigDecimal value,
            java.math.BigDecimal min,
            java.math.BigDecimal max) {
        return value.max(min).min(max);
    }

    protected String substr(String s, int start, int end) {
        int len = s.length();
        start = clamp(start, 0, len);
        end = clamp(end, start, len);
        return s.substring(start, end);
    }

    protected String substr(String s, int start) {
        int len = s.length();
        start = clamp(start, 0, len);
        return s.substring(start, len);
    }

    protected DateTime now() {
        return new DateTime();
    }

    protected LocalDate today() {
        return new LocalDate();
    }

    protected LocalDate to_date(DateTime dt) {
        return dt.toLocalDate();
    }

    private final HashMap<String, Integer> autoincs =
        new HashMap<String, Integer>();

    private final HashMap<String, Pattern> patterns =
        new HashMap<String, Pattern>();
}
