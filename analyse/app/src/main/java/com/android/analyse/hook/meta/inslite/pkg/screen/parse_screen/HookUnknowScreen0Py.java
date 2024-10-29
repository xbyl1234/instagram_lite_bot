package com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookUnknowScreen0Py extends ClassLoadCallBack {
    public HookUnknowScreen0Py(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
//        XposedBridge.hookAllMethods(clz, "A04", new XC_MethodHook() {
//            @Override
//            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                super.beforeHookedMethod(param);
//                String logs = "";
//                logs += "judge: " + frida_helper.byte_2_hex((byte) param.args[0]) + " - " + param.args[1] + ", " +
//                        (((byte) param.args[0] & 1 << (int) param.args[1] % 8) != 0);
//                write_log(logs);
//            }
//        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookUnknowScreen0Py_A03, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                write_log("HookUnknowScreen0Py A03 before: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false));
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                write_log("HookUnknowScreen0Py A03 after: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false));
            }
        });
    }
}
