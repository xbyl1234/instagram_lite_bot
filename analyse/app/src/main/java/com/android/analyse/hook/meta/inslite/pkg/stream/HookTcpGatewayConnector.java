package com.android.analyse.hook.meta.inslite.pkg.stream;

import android.util.Log;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.FieldName;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookTcpGatewayConnector extends ClassLoadCallBack {
    public HookTcpGatewayConnector(AppFileWriter logFile) {
        super(logFile);
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, FieldName.HookTcpGatewayConnectorChangeConnectionStateTo, int.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "changeConnectionStateTo: " + param.args[0] + "\n";
                logs += Log.getStackTraceString(new Throwable());
                write_log(logs);
            }
        });
    }
}
