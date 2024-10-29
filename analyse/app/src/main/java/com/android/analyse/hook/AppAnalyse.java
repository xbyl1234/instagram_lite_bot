package com.android.analyse.hook;

import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.os.UserHandle;
import android.util.Log;

import com.common.log;
import com.common.tools.hooker.FakeClass;
import com.common.tools.hooker.FakeClassBase;
import com.common.tools.hooker.FakeMethod;
import com.common.tools.hooker.Hooker;

import java.io.File;
import java.nio.ByteBuffer;
import java.util.concurrent.ConcurrentHashMap;

import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.ShortBufferException;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XC_MethodReplacement;
import de.robv.android.xposed.XposedBridge;
import de.robv.android.xposed.XposedHelpers;

public class AppAnalyse {
    static public void DoNotDel(ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(File.class, "delete",
                new XC_MethodReplacement() {
                    @Override
                    protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
                        File f = (File) param.thisObject;
                        log.i("delete: " + f.getAbsolutePath());
                        return true;
                    }
                });

        XposedBridge.hookAllConstructors(File.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                File f = (File) param.thisObject;
                log.i("open: " + f.getAbsolutePath());
            }
        });
    }

    static public boolean HookMap() {
        XposedHelpers.findAndHookMethod(ConcurrentHashMap.class, "put", Object.class, Object.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                if (param.args[0].getClass() == String.class) {
                    String key = (String) param.args[0];
                    if (key.length() == 2) {
                        log.i("" + param.args[0]);
                        if (key.equals("tc")) {
                            new Throwable().printStackTrace();
                        }
                    }
                }
            }
        });
        return true;
    }

    static public void HookActivity(ClassLoader classLoader) {
        Class ContextImpl = XposedHelpers.findClass("android.app.ContextImpl", classLoader);
        XposedHelpers.findAndHookMethod(ContextImpl, "startActivity",
                Intent.class, Bundle.class, new XC_MethodHook() {
                    @Override
                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                        super.beforeHookedMethod(param);
                        log.i("startActivity :" + param.args[0] + "\n" + Log.getStackTraceString(new Throwable()));
                    }
                });

        XposedHelpers.findAndHookMethod(ContextImpl, "startActivityAsUser",
                Intent.class, Bundle.class, UserHandle.class, new XC_MethodHook() {
                    @Override
                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                        super.beforeHookedMethod(param);
                        log.i("startActivityAsUser :" + param.args[0] + "\n" + Log.getStackTraceString(new Throwable()));
                    }
                });

        XposedHelpers.findAndHookMethod(ContextImpl, "startActivities", Intent[].class, Bundle.class,
                new XC_MethodHook() {
                    @Override
                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                        super.beforeHookedMethod(param);
                        log.i("startActivities :" + param.args[0] + "\n" + Log.getStackTraceString(new Throwable()));
                    }
                });

        XposedHelpers.findAndHookMethod(ContextImpl, "startActivitiesAsUser", Intent[].class, Bundle.class, UserHandle.class,
                new XC_MethodHook() {
                    @Override
                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                        super.beforeHookedMethod(param);
                        log.i("startActivitiesAsUser :" + param.args[0] + "\n" + Log.getStackTraceString(new Throwable()));
                    }
                });


        XposedHelpers.findAndHookMethod(PendingIntent.class, "getActivityAsUser",
                Context.class, int.class, Intent.class, int.class, Bundle.class, UserHandle.class,
                new XC_MethodHook() {
                    @Override
                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                        super.beforeHookedMethod(param);
                        log.i("getActivityAsUser :" + param.args[0] + "\n" + Log.getStackTraceString(new Throwable()));
                    }
                });

//        Class ApplicationPackageManager = XposedHelpers.findClass("android.app.ApplicationPackageManager", classLoader);
//        XposedHelpers.findAndHookMethod(ApplicationPackageManager, "queryIntentActivitiesAsUser",
//                Intent.class, PackageManager.ResolveInfoFlags.class, int,
//        new XC_MethodHook() {
//            @Override
//            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                super.beforeHookedMethod(param);
//                log.i("getActivityAsUser :" + param.args[0] + "\n" + Log.getStackTraceString(new Throwable()));
//            }
//        });
//
//
//        List<ResolveInfo> queryIntentActivitiesAsUser (userId)


    }

    static public void HookFile(ClassLoader classLoader) {
        XposedBridge.hookAllConstructors(File.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                log.i("new file: " + ((File) param.thisObject).getAbsolutePath());
            }
        });
    }


    static public void HookCipher(ClassLoader classLoader) {
    }

    static public void Hook(ClassLoader classLoader) {
//        HookFile(classLoader);
        HookCipher(classLoader);
    }
}
