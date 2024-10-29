package com.android.analyse.hook.meta.inslite.pkg.screen.view;

import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.common.tools.hooker.HookTools;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;
import de.robv.android.xposed.XposedHelpers;

public class HookSubWrapScreen2 extends ClassLoadCallBack {
    WrapScreenImpl impl;

    public HookSubWrapScreen2(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        impl = new WrapScreenImpl(logFile, "subWrapScreen2", clz, classLoader);
        XposedBridge.hookAllMethods(clz, MethodNames.HookSubWrapScreen2_add_sub_screen, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "add sub screen this: " + param.thisObject + ", sub: " + param.args[0] + "\n";
                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });

        XposedHelpers.findAndHookMethod(clz, MethodNames.HookSubWrapScreen2_on_get_screen_for_resource_id, new XC_MethodHook() {
            int deep = 0;

            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                deep += 1;
                String logs = "before on_get_screen_for_resource_id deep " + deep + " this " + param.thisObject;
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                super.beforeHookedMethod(param);
                deep -= 1;
                String logs = "after on_get_screen_for_resource_id deep " + deep + " this " + param.thisObject + " result: " + param.getResult();
                if (param.getResult() != null) {
                    logs += " resource id:" + HookTools.GetFieldValue(GetClz("X.0Pu"), param.getResult(), "A1M");
                } else {
                    logs += " resource id -0-";
                }
                write_log(logs);
            }
        });


    }
}
