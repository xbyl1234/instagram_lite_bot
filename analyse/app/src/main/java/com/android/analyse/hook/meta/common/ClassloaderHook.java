package com.android.analyse.hook.meta.common;

import com.common.log;

import java.util.HashMap;
import java.util.Map;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class ClassloaderHook {
    Map<String, ClassLoadCallBack> ClassLoadReg = new HashMap<>();
    final Map<String, Class> ClassMaps = new HashMap<>();
    ClassLoader classLoader = null;
    static ClassloaderHook self;

    public ClassloaderHook(String clzName, ClassLoader xpClassLoader) {
        self = this;
        this.classLoader = xpClassLoader;
        XposedHelpers.findAndHookMethod(XposedHelpers.findClass(clzName, xpClassLoader),
                "loadClass", String.class, boolean.class, new XC_MethodHook() {
                    @Override
                    protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                        super.afterHookedMethod(param);
                        String clzName = (String) param.args[0];
                        synchronized (ClassMaps) {
                            ClassMaps.put(clzName, (Class) param.getResult());
                        }
                        if (ClassLoadReg.get(clzName) != null) {
                            log.i("OnLoadedClass " + clzName);
                            ClassLoadReg.get(clzName).BaseOnLoadedClass(clzName, (Class) param.getResult(), classLoader);
                        }
                    }
                });
    }

    public void registerCallback(String clzName, ClassLoadCallBack classLoadCallBack) {
        ClassLoadReg.put(clzName, classLoadCallBack);
    }

    public static Class GetClass(String clzName) {
        synchronized (self.ClassMaps) {
            Class clz = self.ClassMaps.get(clzName);
            if (clz != null) {
                return clz;
            }
        }
        return XposedHelpers.findClass(clzName, self.classLoader);
    }
}
