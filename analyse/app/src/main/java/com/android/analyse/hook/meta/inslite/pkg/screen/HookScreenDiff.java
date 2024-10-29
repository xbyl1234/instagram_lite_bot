package com.android.analyse.hook.meta.inslite.pkg.screen;

import static com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen.dumpViewTree;
import static com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen.tree2string;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookScreenDiff extends ClassLoadCallBack {
    public HookScreenDiff(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllMethods(clz, MethodNames.HookScreenDiff_read_ScreenDiff_item_A03, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "read_ScreenDiff_item_A03 before: ";
                try {
                    DumpScreen.TreeNode tree = dumpViewTree(param.args[2]);
                    logs += "screen obj: " + param.args[2] + "\n";
                    logs += "dump tree1: " + tree2string(tree) + "\n";
                    write_log(logs);
                } catch (Throwable e) {
                    write_log(logs);
                    write_log("error: " + e);
                }
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "read_ScreenDiff_item_A03 after: ";
                try {
                    DumpScreen.TreeNode tree = dumpViewTree(param.args[2]);
                    logs += "screen obj: " + param.args[2] + "\n";
                    logs += "dump tree1: " + tree2string(tree) + "\n";
                    write_log(logs);
                } catch (Throwable e) {
                    write_log(logs);
                    write_log("error: " + e);
                }
            }
        });
    }
}
