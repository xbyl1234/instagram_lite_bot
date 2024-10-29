package com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookSubScreen1 extends ClassLoadCallBack {
    SubScreenImpl impl;

    public HookSubScreen1(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        impl = new SubScreenImpl(logFile, "screen1", clz, classLoader);
        XposedBridge.hookAllMethods(clz, MethodNames.HookSubScreen1_A04, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                write_log("HookSubScreen1 A04 before: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false));
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                write_log("HookSubScreen1 A04 after: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false));
            }
        });
    }
}
