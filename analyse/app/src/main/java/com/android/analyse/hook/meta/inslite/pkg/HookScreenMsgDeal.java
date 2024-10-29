package com.android.analyse.hook.meta.inslite.pkg;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookScreenMsgDeal extends ClassLoadCallBack {
    public HookScreenMsgDeal(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllMethods(clz, MethodNames.HookScreenMsgDeal_call_call_set_data_array_A06, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                write_log("call_call_set_data_array_A06: " + param.args[2] + ", " + Message.Msg2Json(param.args[1], false));
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookScreenMsgDeal_A09, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "";
                logs += "HookScreenMsgDeal A09 before: ";
                logs += param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "";
                logs += "HookScreenMsgDeal A09 after: ";
                logs += param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookScreenMsgDeal_ScreenDecode, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "HookScreenMsgDeal ScreenDecode A0A before: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "HookScreenMsgDeal ScreenDecode A0A after: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookScreenMsgDeal_ScreenDecodeBody, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "HookScreenMsgDeal ScreenDecodeBody A07 before: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "HookScreenMsgDeal ScreenDecodeBody A07 after: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }
        });
        XposedBridge.hookAllMethods(clz, MethodNames.HookScreenMsgDeal_ScreenDataArray, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "HookScreenMsgDeal ScreenDataArray A06 before: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "HookScreenMsgDeal ScreenDataArray A06 after: " + param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }
        });
    }
}
