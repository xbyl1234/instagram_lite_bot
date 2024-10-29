package com.android.analyse.hook.meta.inslite.pkg.stream;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookSendQueue extends ClassLoadCallBack {

    public HookSendQueue(AppFileWriter logFile) {
        super(logFile, true);
    }

    boolean first = true;

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
//        XposedBridge.hookAllMethods(clz, "A00", new XC_MethodReplacement() {
//            @Override
//            protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
//                Object msg = param.args[1];
//                if (Message.GetMsgCode(msg) == 4) {
//                    log.i("pass send log: " + Message.Msg2Json(msg, true));
//                    return null;
//                } else {
//                    return XposedBridge.invokeOriginalMethod(param.method, param.thisObject, param.args);
//                }
//            }
//        });
        XposedHelpers.findAndHookMethod(clz, MethodNames.HookSendQueue_get_from_list, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                if (param.getResult() == null) {
                    return;
                }
                write_log("get from list: " + Message.Msg2Json(param.getResult(), true));
//                if (!first) {
//                    while (true) {
//                        Thread.sleep(1000);
//                        write_log("sleep");
//                    }
//                }
//                first = false;
            }
        });
    }
}