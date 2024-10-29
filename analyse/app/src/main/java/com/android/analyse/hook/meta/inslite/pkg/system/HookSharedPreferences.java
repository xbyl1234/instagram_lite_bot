package com.android.analyse.hook.meta.inslite.pkg.system;

import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.common.tools.hooker.HookTools;
import com.frida.frida_helper;

import java.util.Set;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookSharedPreferences extends ClassLoadCallBack {
    static boolean needStack = true;

    public HookSharedPreferences(AppFileWriter logFile) {
        super(logFile);
    }

    void logRead(XC_MethodHook.MethodHookParam param) {
        String logs = "getSharedPreferences: " + param.args[0] + ", default: " + param.args[1] + ", result: " + param.getResult();
        if (needStack) {
            logs += Log.getStackTraceString(new Throwable());
        }
        write_log(logs);
    }

    void logWrite(XC_MethodHook.MethodHookParam param) {
        String logs = "setSharedPreferences: " + param.args[0] + ", value: " + param.args[1];
        if (needStack) {
            logs += Log.getStackTraceString(new Throwable());
        }
        write_log(logs);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, "getAll", new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "getSharedPreferences getAll: " + frida_helper.object_2_string(param.getResult());
                if (needStack) {
                    logs += Log.getStackTraceString(new Throwable());
                }
                write_log(logs);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getString", String.class, String.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                logRead(param);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getStringSet", String.class, Set.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                logRead(param);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getInt", String.class, int.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                logRead(param);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getLong", String.class, long.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                logRead(param);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getFloat", String.class, float.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                logRead(param);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getBoolean", String.class, boolean.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                logRead(param);
            }
        });
        XposedHelpers.findAndHookMethod(clz, "contains", String.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "getSharedPreferences contains: " + (param.getResult());
                if (needStack) {
                    logs += Log.getStackTraceString(new Throwable());
                }
                write_log(logs);
            }
        });

        Class EditorImpl = HookTools.FindClass("android.app.SharedPreferencesImpl$EditorImpl");
        XposedHelpers.findAndHookMethod(EditorImpl, "putString", String.class, String.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                logWrite(param);
            }
        });


        XposedHelpers.findAndHookMethod(EditorImpl, "putStringSet", String.class, Set.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                logWrite(param);
            }
        });
        XposedHelpers.findAndHookMethod(EditorImpl, "putInt", String.class, int.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                logWrite(param);
            }
        });
        XposedHelpers.findAndHookMethod(EditorImpl, "putLong", String.class, long.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                logWrite(param);
            }
        });
        XposedHelpers.findAndHookMethod(EditorImpl, "putFloat", String.class, float.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                logWrite(param);
            }
        });
        XposedHelpers.findAndHookMethod(EditorImpl, "putBoolean", String.class, boolean.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                logWrite(param);
            }
        });
        XposedHelpers.findAndHookMethod(EditorImpl, "remove", String.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "setSharedPreferences: remove" + param.args[0];
                if (needStack) {
                    logs += Log.getStackTraceString(new Throwable());
                }
                write_log(logs);
            }
        });
        XposedHelpers.findAndHookMethod(EditorImpl, "clear", new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "setSharedPreferences: clear";
                write_log(logs);
            }
        });

    }
}
