package com.android.analyse.hook.meta.inslite.pkg.screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.frida.frida_helper;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookBloksHostingComponent extends ClassLoadCallBack {
    public HookBloksHostingComponent(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllConstructors(clz, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "";
                logs += "screem code: " + param.args[3] + " ";
                logs += "info: " + frida_helper.object_2_string(param.args[1]);
                write_log(logs);
            }
        });
    }
}
