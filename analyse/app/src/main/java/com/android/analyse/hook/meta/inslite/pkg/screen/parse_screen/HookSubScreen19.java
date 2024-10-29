package com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

public class HookSubScreen19 extends ClassLoadCallBack {
    SubScreenImpl impl;

    public HookSubScreen19(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        impl = new SubScreenImpl(logFile, "screen19", clz, classLoader);
    }
}
