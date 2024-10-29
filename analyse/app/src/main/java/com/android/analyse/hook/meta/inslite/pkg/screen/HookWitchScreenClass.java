package com.android.analyse.hook.meta.inslite.pkg.screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;

import java.util.ArrayList;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookWitchScreenClass extends ClassLoadCallBack {
    public HookWitchScreenClass(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, MethodNames.HookWitchScreenClass_witch_screen_class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                write_log("witch screen class: self: " + param.thisObject + " ret: " + param.getResult());
            }
        });
        XposedHelpers.findAndHookMethod(clz, MethodNames.HookWitchScreenClass_get_all_screen_A8P, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "get_all_screen_A8P: ";
                logs += "this: " + param.thisObject + ", ";
                logs += "ret: ";
                ArrayList ret = (ArrayList) param.getResult();
                if (ret != null) {
                    for (int i = 0; i < ret.size(); i++) {
                        logs += ret.get(i) + ", ";
                    }
                }
                write_log(logs);
            }
        });
    }
}
