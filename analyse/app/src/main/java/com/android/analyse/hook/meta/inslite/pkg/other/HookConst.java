package com.android.analyse.hook.meta.inslite.pkg.other;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

import com.frida.Reflect2Json;

public class HookConst extends ClassLoadCallBack {
    public HookConst(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                while (true) {
                    try {
                        write_log("ins const: " + Reflect2Json.StaticClass2Json(clz));
                    } catch (Throwable e) {
                        e.printStackTrace();
                    }
                    try {
                        Thread.sleep(1500);
                    } catch (InterruptedException e) {
                    }
                }
            }
        }).start();

    }
}
