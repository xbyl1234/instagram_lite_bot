package com.android.analyse.hook.meta.inslite.pkg;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

public class HookScreenMsgDeal_0KW extends ClassLoadCallBack {
    public HookScreenMsgDeal_0KW(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
//        XposedBridge.hookAllMethods(clz, "A01", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                String logs = "get_screen_A01: ";
//                logs += "screen: " + param.args[0] + ", ";
//                logs += "param1: " + param.args[1] + ", ";
//                logs += "param2: " + param.args[2] + ", ";
//                logs += "result: " + param.getResult();
//                write_log(logs);
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "A0A", int.class, int.class, new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                String logs = "get_screen_A0A: ";
//                logs += "screenId: " + param.args[0] + ", ";
//                logs += "param1: " + param.args[1] + ", ";
//                logs += "result: " + param.getResult() + ", ";
//                write_log(logs);
//            }
//        });
    }
}
