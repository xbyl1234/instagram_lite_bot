package com.android.analyse.hook.meta.inslite.pkg.screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

import java.util.Map;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookBloksScreen extends ClassLoadCallBack {
    public HookBloksScreen(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookConstructor(clz, String.class, String.class, Map.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "BloksScreen: ";
                logs += "params1: " + param.args[0] + ", ";
                logs += "params2: " + param.args[1] + ", ";
                logs += "params3: " + param.args[2] + ", ";
                write_log(logs);
            }
        });
    }
}
