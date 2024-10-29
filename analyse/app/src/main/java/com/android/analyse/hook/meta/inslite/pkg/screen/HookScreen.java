package com.android.analyse.hook.meta.inslite.pkg.screen;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.MethodNames;
import com.frida.frida_helper;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedBridge;

public class HookScreen extends ClassLoadCallBack {


    public HookScreen(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {

        XposedBridge.hookAllMethods(clz, MethodNames.HookScreen_call_send_action_A00, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
//                0Th param0,
//                index_int_array_map_0Te param1,
//                ArrayList list2,
//                ArrayList list3,
//                List list4,
//                int from_screen_id_5,
//                int to_screen_id_6,
//                int resource_id_7,
//                int param8,
//                short param9,
//                short param10,
//                boolean param11,
//                boolean isCode83_12,
//                boolean param13,
//                boolean param14,
//                boolean cmd_data_has_extern15,
//                boolean param16,
//                boolean param17
                String logs = "deal_screen call_send_action_A00: ";
                logs += "from_screen_id : " + param.args[5] + ", ";
                logs += "to_screen_id : " + param.args[6] + ", ";
                logs += "resource_id : " + param.args[7] + ", ";
                logs += "cmd_data_has_extern : " + param.args[15] + ", ";
                logs += "isCode83_12 : " + param.args[12] + ", ";
                logs += "param0 : " + frida_helper.object_2_string(param.args[0]) + ", ";
                logs += "param8 : " + param.args[8] + ", ";
                logs += "param9 : " + param.args[9] + ", ";
                logs += "param10 : " + param.args[10] + ", ";
                logs += "param11 : " + param.args[11] + ", ";
                logs += "param13 : " + param.args[13] + ", ";
                logs += "param14 : " + param.args[14] + ", ";
                logs += "param16 : " + param.args[16] + ", ";
                logs += "param17 : " + param.args[17] + ", ";
                logs += "list2 : " + frida_helper.object_2_string(param.args[2]) + ", ";
                logs += "list3 : " + frida_helper.object_2_string(param.args[3]) + ", ";
                logs += "list4 : " + frida_helper.object_2_string(param.args[4]) + ", ";
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "after deal_screen call_send_action_A00";
                write_log(logs);
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookScreen_deal_screen_code_A03, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "deal_screen code A03: " + param.args[3] + " ";
                logs += "from: " + param.args[4] + ",  ";
                logs += "data: " + frida_helper.byte_2_hex_str(param.args[2]) + "\n";
//                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "after deal_screen A03";
                write_log(logs);
            }
        });

        XposedBridge.hookAllMethods(clz, MethodNames.HookScreen_deal_screen_code_A02, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "deal_screen code A02: params_code: " + param.args[3] + " ";
                logs += "from: " + param.args[2] + "\n";
//                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }

            @Override
            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
                super.afterHookedMethod(param);
                String logs = "after deal_screen A02";
                write_log(logs);
            }
        });


    }
}
