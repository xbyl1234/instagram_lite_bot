package com.android.analyse.hook.meta.common;

import com.android.analyse.hook.AppFileWriter;
import com.common.log;


abstract public class ClassLoadCallBack {
    public AppFileWriter logFile;
    public boolean adbLog = true;

    public ClassLoadCallBack(AppFileWriter logFile) {
        this.logFile = logFile;
    }

    public ClassLoadCallBack(AppFileWriter logFile, boolean adbLog) {
        this.logFile = logFile;
        this.adbLog = adbLog;
    }

    public abstract void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader);

    public void BaseOnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        OnLoadedClass(clzName, clz, classLoader);
    }

    public void write_log(String data) {
        if (adbLog) {
            log.LogNoWriter(log.LogLevel.Info, data);
        }
        logFile.write(android.os.Process.myTid() + "\t" + data);
    }

    public static Class GetClz(String clzName) {
        return ClassloaderHook.GetClass(clzName);
    }
}