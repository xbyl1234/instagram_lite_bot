package com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

public class HookSubScreen13 extends ClassLoadCallBack {
    SubScreenImpl impl;

    public HookSubScreen13(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        impl = new SubScreenImpl(logFile, "screen13", clz, classLoader);
    }
}
