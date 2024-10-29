package com.android.analyse.hook.meta.inslite.pkg.screen;

import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;
import de.robv.android.xposed.XposedHelpers;

public class HookWitchScreen extends ClassLoadCallBack {
    public HookWitchScreen(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, MethodNames.HookWitchScreen_get_witch_screen_decode, byte.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                write_log("get witch screen decode: " + param.args[0] + ", clz:" + param.getResult());
            }
        });
        XposedBridge.hookAllMethods(clz, MethodNames.HookWitchScreen_create_sub_screen, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "before create sub screen: " + Message.Msg2Json(param.args[0], false);
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "after create sub screen: " + param.getResult() + ", data: " + Message.Msg2Json(param.args[0], false);
                write_log(logs);
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookWitchScreen_read_sub_screen_A03, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "read_sub_screen_A03: ";
                logs += "screen: " + param.args[0] + ", ";
                logs += "type: " + param.args[2] + ", ";
                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });
    }
}
