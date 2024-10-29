package com.android.analyse.hook;


import android.app.Application;

import com.android.analyse.hook.meta.fblite.HookFbLite;
import com.android.analyse.hook.meta.inslite.HookInsLite;
import com.android.analyse.hook.meta.sslpinning.PassSsl;
import com.common.log;
import com.common.tools.hooker.WhenHook;

import java.io.File;

import de.robv.android.xposed.IXposedHookLoadPackage;
import de.robv.android.xposed.callbacks.XC_LoadPackage;

public class AnalyseHook implements IXposedHookLoadPackage {
    static boolean hadHook = false;

    @Override
    public void handleLoadPackage(XC_LoadPackage.LoadPackageParam lpparam) throws Throwable {
        if (lpparam.processName.contains("com.android.analyse") || lpparam.processName.contains("com.android.webview") || hadHook) {
            return;
        }
        hadHook = true;
        log.i("analyse inject process: " + lpparam.processName);
        try {
            new File("/sdcard/Android/data/" + lpparam.packageName).mkdir();
        } catch (Throwable e) {
            log.i("make extern dir error: " + e);
        }
        FridaHelperLoader.InjectFridaHelp(lpparam.classLoader);

        if (lpparam.packageName.equals("android")) {
            log.i("will inject");
        } else {
            Native.LoadAnalyseSo(lpparam.packageName);
            WhenHook.WhenPerformLaunchActivityHook(new WhenHook.WhenHookCallback() {
                @Override
                public void OnHook(Application application, ClassLoader classLoader) {
                    try {
                        application.getExternalCacheDir();
                    } catch (Throwable e) {
                        e.printStackTrace();
                    }
                }
            });
            if (lpparam.packageName.equals("com.instagram.lite")) {
                HookInsLite.Hook(lpparam.classLoader);
            }
            if (lpparam.packageName.equals("com.facebook.lite")) {
                HookFbLite.Hook(lpparam.classLoader);
            }
            if (lpparam.packageName.equals("com.facebook.katana")) {
                PassSsl.HookFb(lpparam.classLoader);
            }
            if (lpparam.packageName.equals("com.instagram.android")) {
                PassSsl.HookIns(lpparam.classLoader);
            }
        }

    }
}
