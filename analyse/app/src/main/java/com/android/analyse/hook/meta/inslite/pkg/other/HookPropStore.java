package com.android.analyse.hook.meta.inslite.pkg.other;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.frida.frida_helper;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookPropStore extends ClassLoadCallBack {
    public HookPropStore(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, MethodNames.HookPropStore_set, int.class, Object.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "HookPropStore: ";
                logs += "idx: " + param.args[0] + ", ";
                logs += "data: " + frida_helper.object_2_string(param.args[1]) + ", " + "\n";
//                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });

    }
}
