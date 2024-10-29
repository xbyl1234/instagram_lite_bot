package com.android.analyse.hook.meta.inslite.pkg.stream;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.common.log;

import java.io.InputStream;
import java.io.OutputStream;

import de.robv.android.xposed.XC_MethodReplacement;
import de.robv.android.xposed.XposedBridge;
import de.robv.android.xposed.XposedHelpers;

public class JavaSocketWrapHooker extends ClassLoadCallBack {
    public JavaSocketWrapHooker(AppFileWriter logFile) {
        super(logFile,false);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, "getInputStream", new XC_MethodReplacement() {
            @Override
            protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
                InputStream in = (InputStream) XposedBridge.invokeOriginalMethod(param.method, param.thisObject, param.args);
                log.i("create MyInputStream");
                return new MyInputStream(in, logFile, "tcp");
            }
        });
        XposedHelpers.findAndHookMethod(clz, "getOutputStream", new XC_MethodReplacement() {
            @Override
            protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
                OutputStream out = (OutputStream) XposedBridge.invokeOriginalMethod(param.method, param.thisObject, param.args);
                log.i("create MyOutputStream");
                return new MyOutputStream(out, logFile);
            }
        });
    }
}
