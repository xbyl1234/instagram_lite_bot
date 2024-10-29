package com.android.analyse.hook.meta.inslite.pkg.screen.view;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

public class HookSubWrapScreen1 extends ClassLoadCallBack {
    WrapScreenImpl impl;
    public HookSubWrapScreen1(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        impl = new WrapScreenImpl(logFile, "subWrapScreen1", clz, classLoader);
    }

}
