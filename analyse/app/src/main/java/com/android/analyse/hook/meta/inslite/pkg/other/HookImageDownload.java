package com.android.analyse.hook.meta.inslite.pkg.other;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.msg.Message;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookImageDownload extends ClassLoadCallBack {
    public HookImageDownload(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllMethods(clz, MethodNames.HookImageDownload_image_download, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "image download: ";
                logs += "buff: " + Message.MsgOffset(param.args[1]) + ", ";
                logs += "id: " + param.args[6] + ", ";
                logs += "totalSize: " + param.args[2] + ", ";
                logs += "offset: " + param.args[3] + ", ";
                logs += "dataLength: " + param.args[4] + ", ";
                logs += " flags: " + param.args[5] + ", ";
                write_log(logs);
            }
        });
    }
}
