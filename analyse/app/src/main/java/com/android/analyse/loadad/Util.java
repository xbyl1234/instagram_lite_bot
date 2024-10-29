package com.android.analyse.loadad;

import com.fucker.gson.Gson;

import java.lang.reflect.Type;

public class Util
{

    public static String toJson(Object obj) {
        return new Gson().toJson(obj);
    }

    public static <T> T fromJson(String text, Class<T> clazz) {
        return new Gson().fromJson(text, clazz);
    }

    public static <T> T fromJsonType(String text, Type type) {
        return new Gson().fromJson(text, type);
    }

}
