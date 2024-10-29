package com.android.analyse.hook.meta.inslite.pkg.system;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookArrayList extends ClassLoadCallBack {
    public HookArrayList(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, "get", int.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                Object obj = param.getResult();
                if (obj != null) {
                    if (obj.getClass().getName().contains("X.")) {
                        String logs = "list get: " + param.args[0] + ": " + obj;
                        write_log(logs);
                    }
                }
            }
        });
    }
}
