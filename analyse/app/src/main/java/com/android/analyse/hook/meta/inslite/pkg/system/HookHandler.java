package com.android.analyse.hook.meta.inslite.pkg.system;

import android.os.Message;
import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.pkg.other.HookWindowsMsg;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookHandler extends ClassLoadCallBack {
    public HookHandler(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, "sendMessage", android.os.Message.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                Message msg = (Message) param.args[0];
                if (msg.obj != null && msg.obj.getClass().getName().equals("X.0K0")) {
                    String logs = "sendMessage: ";
                    logs += "msg: " + HookWindowsMsg.MsgObj2Str(msg.obj);
                    logs += Log.getStackTraceString(new Throwable());
                    write_log(logs);
                }
            }
        });
        XposedHelpers.findAndHookMethod(clz, "sendMessageDelayed", android.os.Message.class, long.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                Message msg = (Message) param.args[0];
                if (msg.obj != null && msg.obj.getClass().getName().equals("X.0K0")) {
                    String logs = "sendMessageDelayed: ";
                    logs += "msg: " + HookWindowsMsg.MsgObj2Str(msg.obj);
                    logs += Log.getStackTraceString(new Throwable());
                    write_log(logs);
                }
            }
        });
    }
}
