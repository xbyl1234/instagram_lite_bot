package com.android.analyse.hook.meta.inslite.pkg.msg;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;

public class SlicedByteBufferHooker extends ClassLoadCallBack {

    public SlicedByteBufferHooker(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {

    }

    public void HookRecvRead(Class clz) {
//        XposedHelpers.findAndHookMethod(clz, "ATn", int.class, new XC_MethodHook() {
//            @Override
//            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                super.beforeHookedMethod(param);
//                String logs = "";
//                logs += "set offset:";
//                logs += frida_helper.object_2_string(param.args[0]) + ", ";
//                logs += frida_helper.object_2_string(param.thisObject) + "\n";
//                logs += Log.getStackTraceString(new Throwable());
//                write_log(logs);
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "A7u", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readBoolean: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "A81", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readByte: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "readDouble", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readDouble: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "readFloat", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readFloat: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "readFully", byte[].class, int.class, int.class, new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "AA9", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readInt: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "AAm", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "--readLong: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "ADY", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readShort: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "AEQ", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readUnsignedByte: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "AER", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "readUnsignedShort: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "ADa", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "read_signed_varint32_int: " + param.getResult());
//            }
//        });
//        XposedHelpers.findAndHookMethod(clz, "ABW", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "read_str_array_varlen: " + param.getResult());
//            }
//        });
//
//        XposedHelpers.findAndHookMethod(clz, "ABV", new XC_MethodHook() {
//            @Override
//            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                super.afterHookedMethod(param);
//                logInfo.PushLogs(param.thisObject, "read_string_varlen: " + param.getResult());
//            }
//        });
    }

}
