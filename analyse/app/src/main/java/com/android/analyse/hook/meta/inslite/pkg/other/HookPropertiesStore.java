package com.android.analyse.hook.meta.inslite.pkg.other;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookPropertiesStore extends ClassLoadCallBack {
    public HookPropertiesStore(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, MethodNames.HookPropertiesStore_get_int_key, int.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                write_log("PropertiesStore: " + param.thisObject + " get int key: " + param.args[0] + " value: " + param.getResult());
            }
        });
    }
}
