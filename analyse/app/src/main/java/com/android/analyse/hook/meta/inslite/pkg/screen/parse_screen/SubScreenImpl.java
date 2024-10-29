package com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;
import com.common.log;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class SubScreenImpl extends ClassLoadCallBack {
    volatile static int count = 0;

    public SubScreenImpl(AppFileWriter logFile, String screenName, Class clz, ClassLoader classLoader) {
        super(logFile);

        XposedBridge.hookAllMethods(clz, MethodNames.SubScreenImpl_A09, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                write_log(screenName + " A09 before: " + param.args[1] + " " + Message.Msg2Json(param.args[1],false));
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                write_log(screenName + " A09 after: " + param.args[1] + " " + Message.Msg2Json(param.args[1],false));
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.SubScreenImpl_A0B, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                synchronized (this) {
                    count += 1;
                    String logs = screenName + " before count " + count + " screen2 read_ScreenChange_A0B: ";
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
                    String logs = screenName + " after count " + count + " read_ScreenChange_A0B: ";
                    write_log(logs);
                    count -= 1;
                }
            }
        });
        log.i("hook " + screenName + " success");
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {

    }
}
