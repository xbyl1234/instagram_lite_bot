package com.android.analyse.hook.meta.fblite;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassloaderHook;
import com.android.analyse.hook.meta.inslite.pkg.stream.HookRecvQueue;
import com.android.analyse.hook.meta.inslite.pkg.stream.HookSendQueue;
import com.android.analyse.hook.meta.inslite.pkg.stream.JavaSocketWrapHooker;

public class HookFbLite {
    static ClassloaderHook classloaderHook;
    static public AppFileWriter pkgSendRecvFile = new AppFileWriter("sr");
    static public AppFileWriter streamFile = new AppFileWriter("stream");
    static public AppFileWriter analyse = new AppFileWriter("analyse");

    public static void Hook(ClassLoader classLoader) {
        if (classloaderHook != null) {
            return;
        }
        ClassNames.init_401_0_0_14_110_503500325();
        classloaderHook = new ClassloaderHook(com.android.analyse.hook.meta.inslite.ClassNames.loadClass, classLoader);

        classloaderHook.registerCallback(ClassNames.javaSocketWrap, new JavaSocketWrapHooker(streamFile));
        classloaderHook.registerCallback(ClassNames.sendQueue, new HookSendQueue(pkgSendRecvFile));
        classloaderHook.registerCallback(ClassNames.recvQueue, new HookRecvQueue(pkgSendRecvFile));

    }
}
