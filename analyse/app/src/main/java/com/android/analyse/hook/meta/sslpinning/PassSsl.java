package com.android.analyse.hook.meta.sslpinning;

import com.android.analyse.hook.Native;

import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.ArrayList;

import javax.net.ssl.KeyManager;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XC_MethodReplacement;
import de.robv.android.xposed.XposedBridge;
import de.robv.android.xposed.XposedHelpers;

public class PassSsl {
    static class MyX509TrustManager implements X509TrustManager {
        @Override
        public void checkClientTrusted(X509Certificate[] chain, String authType) throws CertificateException {

        }

        @Override
        public void checkServerTrusted(X509Certificate[] chain, String authType) throws CertificateException {

        }

        @Override
        public X509Certificate[] getAcceptedIssuers() {
            return new X509Certificate[0];
        }
    }

    static public void HookX509(ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(SSLContext.class, "init",
                KeyManager[].class, TrustManager[].class, SecureRandom.class, new XC_MethodReplacement() {
                    @Override
                    protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
                        MyX509TrustManager[] TrustManager = new MyX509TrustManager[1];
                        TrustManager[0] = new MyX509TrustManager();
                        param.args[1] = TrustManager;
                        return XposedBridge.invokeOriginalMethod(param.method, param.thisObject, param.args);
                    }
                });
        XposedBridge.hookAllMethods(XposedHelpers.findClass("com.android.org.conscrypt.TrustManagerImpl", classLoader), "checkTrustedRecursive",
                new XC_MethodReplacement() {
                    @Override
                    protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
                        return new ArrayList<>();
                    }
                });
    }

    static public void HookIns(ClassLoader classLoader) {
        HookX509(classLoader);
        Native.passInsSslPinning();
    }

    static public void HookFb(ClassLoader classLoader) {
        HookX509(classLoader);
        Native.passFbSslPinning();
    }
}
