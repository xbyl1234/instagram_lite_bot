package com.android.analyse.hook.meta.inslite.pkg.stream;


import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookRecvQueue extends ClassLoadCallBack {

    public HookRecvQueue(AppFileWriter logFile) {
        super(logFile, true);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
//        XposedBridge.hookAllConstructors(RecvQueueClass, new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                log.i("create recv MyQueue");
//                Reflect.SetField(RecvQueueClass, param.thisObject, "A03", new MyQueue());
//                Object A04 = Reflect.GetFieldValue(RecvQueueClass, param.thisObject, "A04");
//                Object A00 = Reflect.GetFieldValue(RecvQueueClass, param.thisObject, "A00");
//                String logs = "";
//                if (A04 != null) {
//                    logs += "A04 is:" + A04.getClass() + " " + A04 + ", ";
//                }
//                if (A00 != null) {
//                    logs += "A00 is:" + A00.getClass() + " " + A00;
//                }
//                log.i(logs);
//                pkgRecvFile.write(logs);
//            }
//        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookRecvQueue_on_recv_msg, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "on recv msg : " + Message.Msg2Json(param.args[1], true);
                write_log(logs);
            }
        });
    }
}
