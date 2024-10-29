package com.android.analyse.hook.meta.inslite.pkg.screen.view;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.frida.frida_helper;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class WrapScreenImpl extends ClassLoadCallBack {
    volatile static int count = 0;

    public WrapScreenImpl(AppFileWriter logFile, String screenName, Class clz, ClassLoader classLoader) {
        super(logFile);
        XposedBridge.hookAllMethods(clz, "A2E", new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "";
                synchronized ("WrapScreenImpl") {
                    logs += screenName + " " + param.thisObject + " set submit data " + count + " before: " +
                            frida_helper.object_2_string(param.args[0]) +
                            ", data: " +
                            frida_helper.object_2_string(param.args[1]) +
                            ", flag: " +
                            param.args[2];
                    count++;
                }
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "";
                synchronized ("WrapScreenImpl") {
                    logs += screenName + " " + param.thisObject + " set submit data " + count + " after: " +
                            frida_helper.object_2_string(param.args[0]) +
                            ", data: " +
                            frida_helper.object_2_string(param.args[1]) +
                            ", flag: " +
                            param.args[2];
                    count--;
                }
                write_log(logs);
            }
        });
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {

    }
}
