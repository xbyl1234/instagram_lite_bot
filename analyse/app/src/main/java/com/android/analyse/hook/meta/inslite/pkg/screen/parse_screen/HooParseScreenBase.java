package com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HooParseScreenBase extends ClassLoadCallBack {
    public HooParseScreenBase(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllMethods(clz, MethodNames.HooParseScreenBase_A03, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "";
                logs += "HooParseScreenBase A03 before: ";
                logs += param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "";
                logs += "HooParseScreenBase A03 after: ";
                logs += param.args[1] + " " + Message.Msg2Json(param.args[1], false);
                write_log(logs);
            }
        });
        XposedBridge.hookAllMethods(clz, MethodNames.HooParseScreenBase_read_ScreenChange_A0B, new XC_MethodHook() {
            int count = 0;

            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                synchronized (this) {
                    count += 1;
                    String logs = "before count " + count + " screen base read_ScreenChange_A0B: ";
                    logs += "this: " + param.thisObject + ", ";
                    logs += "screen: " + param.args[0] + ", ";
                    logs += "chType: " + param.args[2] + ", ";
                    logs += "data: " + Message.MsgOffset(param.args[1]) + ", ";
//                    logs += Log.getStackTraceString(new Throwable());
                    write_log(logs);
                }
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                synchronized (this) {
                    String logs = "after count " + count + " read_ScreenChange_A0B: ";
                    write_log(logs);
                    count -= 1;
                }
            }
        });
    }
}
