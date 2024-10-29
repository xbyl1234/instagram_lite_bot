package com.android.analyse.hook.meta.inslite.pkg.other;

import static com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen.dumpViewTree;
import static com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen.getScreenId;
import static com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen.getScreenName;
import static com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen.tree2string;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.android.analyse.hook.meta.inslite.pkg.screen.DumpScreen;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookWindowManager extends ClassLoadCallBack {
    public HookWindowManager(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllMethods(clz, MethodNames.HookWindowManager_call_changeScreenTo, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "call_changeScreenTo A1J: ";
                try {
                    logs += "screen1: " + param.args[0] + ", ";
                    logs += "screen2: " + param.args[1] + ", ";
                    logs += "navigationData: " + param.args[2] + ", ";
                    logs += "UnknowFlag: " + param.args[3] + ", ";
                    logs += "receivedScreenId: " + param.args[4] + ", ";
                    logs += "recvInt3: " + param.args[5] + ", ";
                    logs += "bool1: " + param.args[6] + ", ";
                    logs += "bool2: " + param.args[7] + ", ";
                    logs += "displayNow: " + param.args[8] + ", ";
                    logs += "bool4: " + param.args[9] + ", " + "\n";

                    DumpScreen.TreeNode tree = dumpViewTree(param.args[0]);
                    logs += "screen obj: " + param.args[2] + ", id: " + getScreenId(param.args[0]) + ", name: " + getScreenName(param.args[0]) + "\n";
                    logs += "dump tree1: " + tree2string(tree) + "\n";

                    write_log(logs);
                } catch (Throwable e) {
                    write_log(logs);
                    write_log("error: " + e);
                }
            }
        });
        XposedBridge.hookAllMethods(clz, MethodNames.HookWindowManager_A0D, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "WindowManager A0D: ";
                logs += "screen: " + param.args[0] + ", ";
                logs += "screen id: " + param.args[1] + ", ";
                logs += "if show: " + param.args[2] + ", ";
                write_log(logs);
            }
        });
    }
}
