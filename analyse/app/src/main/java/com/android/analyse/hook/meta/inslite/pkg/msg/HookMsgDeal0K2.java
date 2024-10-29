package com.android.analyse.hook.meta.inslite.pkg.msg;

import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.frida.frida_helper;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookMsgDeal0K2 extends ClassLoadCallBack {
    public HookMsgDeal0K2(AppFileWriter logFile) {
        super(logFile, true);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedBridge.hookAllMethods(clz, MethodNames.HookMsgDeal0K2_parse_msg45, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "";
                logs += "parse msg 45:";
                logs += frida_helper.object_2_string(param.args[0]) + "\n";
                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });
        XposedBridge.hookAllMethods(clz, MethodNames.HookMsgDeal0K2_deal_msg, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "";
                logs += "deal msg:";
                logs += frida_helper.object_2_string(param.args[0]) + "\n";
                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });
        XposedBridge.hookAllMethods(clz, MethodNames.HookMsgDeal0K2_like_send_action_A0d, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "like_send_action_A0d: ";
                logs += "param0: " + frida_helper.object_2_string(param.args[0]) + ", ";
                logs += "msg_code_map: " + frida_helper.object_2_string(param.args[1]) + ", ";
                logs += "msg_obj: " + frida_helper.object_2_string(param.args[2]) + ", ";
                logs += "list3: " + frida_helper.object_2_string(param.args[3]) + ", ";
                logs += "from_ScreenId_args_5: " + param.args[4] + ", ";
                logs += "to_ScreenId: " + param.args[5] + ", ";
                logs += "resource_id: " + param.args[6] + ", ";
                logs += "like_action_resource_id_unknow4: " + param.args[7] + ", ";
                logs += "unknow8: " + param.args[8] + ", ";
                logs += "unknow9: " + param.args[9] + ", ";
                logs += "isCode83: " + param.args[10] + ", ";
                logs += "cmd_data_has_extern: " + param.args[11] + ", ";
                logs += "param12: " + param.args[12] + ", ";
                logs += "\n" + Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });
    }
}
