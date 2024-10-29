package com.android.analyse.hook.meta.inslite.pkg.msg;

import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.FieldName;
import com.common.tools.hooker.HookTools;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookSendMsgData extends ClassLoadCallBack {

    public HookSendMsgData(AppFileWriter logFile) {
        super(logFile, true);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookConstructor(clz, int.class, int.class, new XC_MethodHook() {
            @Override
            protected void afterHookedMethod(XC_MethodHook.MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logStr = "";
                logStr = "msg code " + HookTools.GetFieldValue(clz, param.thisObject, FieldName.HookSendMsgData_msg_code) + " create at: \n";
                logStr += Log.getStackTraceString(new Throwable());
                write_log(logStr);
            }
        });
//        XposedHelpers.findAndHookMethod(clz, make_sender_data_to, OutputStream.class, boolean.class, new XC_MethodHook() {
//            @Override
//            protected void beforeHookedMethod(XC_MethodHook.MethodHookParam param) throws Throwable {
//                super.beforeHookedMethod(param);
//                PkgLogHelper.LogMsgData("make_sender_data_to" + param.args[1] + " ", param.thisObject);
//            }
//        });
    }
}
