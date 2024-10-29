package com.android.analyse.hook.meta.inslite;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassloaderHook;
import com.android.analyse.hook.meta.inslite.pkg.other.HookPropStore;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HooParseScreenBase;
import com.android.analyse.hook.meta.inslite.pkg.screen.HookBloksScreen;
import com.android.analyse.hook.meta.inslite.pkg.stream.HookTcpGatewayConnector;
import com.android.analyse.hook.meta.inslite.pkg.system.HookHandler;
import com.android.analyse.hook.meta.inslite.pkg.other.HookImageDownload;
import com.android.analyse.hook.meta.inslite.pkg.msg.HookMsgDeal0K2;
import com.android.analyse.hook.meta.inslite.pkg.stream.HookRecvQueue;
import com.android.analyse.hook.meta.inslite.pkg.screen.HookScreen;
import com.android.analyse.hook.meta.inslite.pkg.screen.HookScreenDiff;
import com.android.analyse.hook.meta.inslite.pkg.HookScreenMsgDeal;
import com.android.analyse.hook.meta.inslite.pkg.HookScreenMsgDeal_0KW;
import com.android.analyse.hook.meta.inslite.pkg.msg.HookSendMsgData;
import com.android.analyse.hook.meta.inslite.pkg.stream.HookSendQueue;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookSubScreen1;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookSubScreen13;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookSubScreen19;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookSubScreen2;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookSubScreen3;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookSubScreen9;
import com.android.analyse.hook.meta.inslite.pkg.screen.parse_screen.HookUnknowScreen0Py;
import com.android.analyse.hook.meta.inslite.pkg.other.HookWindowManager;
import com.android.analyse.hook.meta.inslite.pkg.other.HookWindowsMsg;
import com.android.analyse.hook.meta.inslite.pkg.screen.HookWitchScreen;
import com.android.analyse.hook.meta.inslite.pkg.stream.JavaSocketWrapHooker;
import com.android.analyse.hook.meta.inslite.pkg.msg.SlicedByteBufferHooker;
import com.android.analyse.hook.meta.inslite.pkg.screen.view.HookSubWrapScreen1;
import com.android.analyse.hook.meta.inslite.pkg.screen.view.HookSubWrapScreen13;
import com.android.analyse.hook.meta.inslite.pkg.screen.view.HookSubWrapScreen19;
import com.android.analyse.hook.meta.inslite.pkg.screen.view.HookSubWrapScreen2;
import com.android.analyse.hook.meta.inslite.pkg.screen.view.HookSubWrapScreen3;
import com.android.analyse.hook.meta.inslite.pkg.screen.view.HookSubWrapScreen9;
import com.android.analyse.hook.meta.inslite.pkg.system.HookSharedPreferences;
import com.common.log;

public class HookInsLite {
    static ClassloaderHook classloaderHook;
    static public AppFileWriter pkgSendRecvFile = new AppFileWriter("sr");
    static public AppFileWriter streamFile = new AppFileWriter("stream");
    static public AppFileWriter analyse = new AppFileWriter("analyse");

    public static void Hook(ClassLoader classLoader) {
        if (classloaderHook != null) {
            return;
        }
        ClassNames.init_382_0_0_11_115_538547752();
        MethodNames.init_382_0_0_11_115_538547752();
        FieldName.init_382_0_0_11_115_538547752();
        try {
            classloaderHook = new ClassloaderHook(ClassNames.loadClass, classLoader);

            classloaderHook.registerCallback(ClassNames.javaSocketWrap, new JavaSocketWrapHooker(streamFile));
            classloaderHook.registerCallback(ClassNames.sendQueue, new HookSendQueue(pkgSendRecvFile));
            classloaderHook.registerCallback(ClassNames.recvQueue, new HookRecvQueue(pkgSendRecvFile));
            classloaderHook.registerCallback(ClassNames.TcpGatewayConnector, new HookTcpGatewayConnector(pkgSendRecvFile));

//        classloaderHook.registerCallback(ClassNames.bloksHostingComponent, new HookBloksHostingComponent(analyse));
            //read or write log
//        classloaderHook.registerCallback(ClassNames.slicedByteBuffer, new SlicedByteBufferHooker(analyse));
            //sned create at
            classloaderHook.registerCallback(ClassNames.sendMsgData, new HookSendMsgData(analyse));
            classloaderHook.registerCallback(ClassNames.msgDeal0K2, new HookMsgDeal0K2(analyse));
            classloaderHook.registerCallback(ClassNames.slicedByteBuffer, new SlicedByteBufferHooker(analyse));
            classloaderHook.registerCallback(ClassNames.screenClass, new HookScreen(analyse));
            classloaderHook.registerCallback(ClassNames.screenMsgDeal, new HookScreenMsgDeal(analyse));
            classloaderHook.registerCallback(ClassNames.parseScreenBase, new HooParseScreenBase(analyse));
            classloaderHook.registerCallback(ClassNames.witchScreen, new HookWitchScreen(analyse));
            classloaderHook.registerCallback(ClassNames.unknowScreen0Py, new HookUnknowScreen0Py(analyse));
            classloaderHook.registerCallback(ClassNames.subScreen1, new HookSubScreen1(analyse));
            classloaderHook.registerCallback(ClassNames.subScreen2, new HookSubScreen2(analyse));
            classloaderHook.registerCallback(ClassNames.subScreen3, new HookSubScreen3(analyse));
            classloaderHook.registerCallback(ClassNames.subScreen9, new HookSubScreen9(analyse));
            classloaderHook.registerCallback(ClassNames.subScreen13, new HookSubScreen13(analyse));
            classloaderHook.registerCallback(ClassNames.subScreen19, new HookSubScreen19(analyse));
            classloaderHook.registerCallback(ClassNames.subWrapScreen1, new HookSubWrapScreen1(analyse));
            classloaderHook.registerCallback(ClassNames.subWrapScreen2, new HookSubWrapScreen2(analyse));
            classloaderHook.registerCallback(ClassNames.subWrapScreen3, new HookSubWrapScreen3(analyse));
            classloaderHook.registerCallback(ClassNames.subWrapScreen9, new HookSubWrapScreen9(analyse));
            classloaderHook.registerCallback(ClassNames.subWrapScreen13, new HookSubWrapScreen13(analyse));
            classloaderHook.registerCallback(ClassNames.subWrapScreen19, new HookSubWrapScreen19(analyse));
//            classloaderHook.registerCallback(ClassNames.witchScreenClass, new HookWitchScreenClass(analyse));
            classloaderHook.registerCallback(ClassNames.windowsMsg, new HookWindowsMsg(analyse));
            classloaderHook.registerCallback(ClassNames.bloksScreen, new HookBloksScreen(analyse));
            classloaderHook.registerCallback(ClassNames.propStore, new HookPropStore(analyse));
//            classloaderHook.registerCallback(ClassNames.propertiesStore, new HookPropertiesStore(analyse));
            classloaderHook.registerCallback(ClassNames.screenMsgDeal_0KW, new HookScreenMsgDeal_0KW(analyse));
            classloaderHook.registerCallback(ClassNames.screenDiff, new HookScreenDiff(analyse));
            classloaderHook.registerCallback(ClassNames.windowManager, new HookWindowManager(analyse));
            classloaderHook.registerCallback(ClassNames.imageDownload, new HookImageDownload(analyse));
            classloaderHook.registerCallback(ClassNames.imageDownload, new HookImageDownload(analyse));


//            new HookArrayList(analyse).OnLoadedClass("java.util.ArrayList", Class.forName("java.util.ArrayList"), classLoader);
            new HookHandler(analyse).OnLoadedClass("android.os.Handler", Class.forName("android.os.Handler"), classLoader);
            new HookSharedPreferences(analyse).OnLoadedClass("android.app.SharedPreferencesImpl", Class.forName("android.app.SharedPreferencesImpl"), classLoader);
//            new HookConst(analyse).OnLoadedClass(ClassNames.Const, XposedHelpers.findClass(ClassNames.Const, classLoader), classLoader);
        } catch (Throwable e) {
            log.e(e.toString());
            e.printStackTrace();
        }
    }
}


// DeflaterOutputStream
//        ClassLoadReg.put(className.DeflaterOutputStreamClass, new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                XposedHelpers.findAndHookConstructor(XposedHelpers.findClass(className.DeflaterOutputStreamClass, classLoader),
//                        OutputStream.class, Deflater.class, new XC_MethodHook() {
//                            @Override
//                            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                                super.beforeHookedMethod(param);
//                                log.i("create MyDeflater");
//                                param.args[1] = new MyDeflater2(9, msgDefFile);
//                            }
//                        });
//            }
//        });
//        ClassLoadReg.put("X.0mi", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class msg_header_body = XposedHelpers.findClass(clz, classLoader);
//                XposedHelpers.findAndHookConstructor(msg_header_body, Integer.class, int.class, int.class, new XC_MethodHook() {
//                    @Override
//                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                        super.beforeHookedMethod(param);
//                        String logs = "create msg_header_body: " +
//                                " msg_data_len: " + param.args[1] + ", " +
//                                " is lamz2: " + param.args[0] + ", " +
//                                " stream_idx: " + param.args[2];
//                        log.i(logs);
//                        pkgRecvFile.write(logs);
//                    }
//                });
//            }
//        });
//        ClassLoadReg.put("X.0ml", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class recv_lzma2_stream = XposedHelpers.findClass("X.0ml", classLoader);
//                XposedHelpers.findAndHookMethod(recv_lzma2_stream, "read", byte[].class, int.class, int.class, new XC_MethodHook() {
//                    @Override
//                    protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                        super.afterHookedMethod(param);
//                        String logs = "recv_lzma2_stream decode data: " + frida_helper.byte_2_hex_str(param.args[0]);
//                        log.i(logs);
//                        pkgRecvFile.write(logs);
//                    }
//                });
//                XposedHelpers.findAndHookConstructor(recv_lzma2_stream, InputStream.class, int.class, new XC_MethodReplacement() {
//                    @Override
//                    protected Object replaceHookedMethod(MethodHookParam param) throws Throwable {
//                        String logs = "create recv_lzma2_stream: " + param.args[0] +
//                                " dict size: " + param.args[1] + ", " +
//                                " data: " + frida_helper.byte_2_hex_str((byte[]) Reflect.GetFieldValue(param.args[0].getClass(), param.args[0], "A02"),
//                                0, (int) Reflect.GetFieldValue(param.args[0].getClass(), param.args[0], "A00"));
//                        log.i(logs);
//                        pkgRecvFile.write(logs);
//                        return XposedBridge.invokeOriginalMethod(param.method, param.thisObject, param.args);
//                    }
//                });
//            }
//        });

//        ClassLoadReg.put("X.0aR", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class lzma_input = XposedHelpers.findClass(clz, classLoader);
//                XposedHelpers.findAndHookMethod(lzma_input, "read",
//                        new XC_MethodHook() {
//                            @Override
//                            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                                super.afterHookedMethod(param);
//                                byte tmp[] = new byte[1];
//                                tmp[0] = (byte) ((byte) ((Integer) param.getResult()).intValue() & 0xff);
//                                String logs = "lzma input: " + frida_helper.byte_2_hex_str(tmp);
////                                log.i(logs);
//                                pkgRecvFile.write(logs);
//                            }
//                        });
//
//                XposedHelpers.findAndHookMethod(lzma_input, "read", byte[].class, int.class, int.class,
//                        new XC_MethodHook() {
//                            @Override
//                            protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                                super.afterHookedMethod(param);
//                                String logs = "lzma input: " + frida_helper.byte_2_hex_str(param.args[0], 0, (Integer) param.getResult());
////                                log.i(logs);
//                                pkgRecvFile.write(logs);
//                            }
//                        });
//
//            }
//        });


//        ClassLoadReg.put("X.0K9", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class msgCodeDealFuncWrap = XposedHelpers.findClass(clz, classLoader);
//                XposedBridge.hookAllMethods(msgCodeDealFuncWrap, "A00", new XC_MethodHook() {
//                    @Override
//                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                        super.beforeHookedMethod(param);
//                        String logs = "put msg handle:" + frida_helper.object_2_string(param.args[0]) + "\n";
//                        logs += "all: " + frida_helper.object_2_string(Reflect.GetFieldValue(msgCodeDealFuncWrap, param.thisObject, "A00"));
//                        log.i(logs);
//                        pkgRecvFile.write(logs);
//                    }
//                });
//            }
//        });

//        ClassLoadReg.put("X.0K2", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class msg_deal = XposedHelpers.findClass(clz, classLoader);
//                XposedBridge.hookAllMethods(msg_deal, "AJN",
//                        new XC_MethodHook() {
//                            @Override
//                            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                                super.beforeHookedMethod(param);
//                                Method getMsgCode = XposedHelpers.findClass("X.0Ig", classLoader).getDeclaredMethod("AB3");
//                                int msgCode = (int) getMsgCode.invoke(param.args[0]);
//
//                                Object warpDealFunc = Reflect.GetFieldValue(msg_deal, param.thisObject, "A0Z");
//                                Object dealFunc = Reflect.GetFieldValue(XposedHelpers.findClass("X.0K9", classLoader), warpDealFunc, "A00");
//
//                                Method get = dealFunc.getClass().getDeclaredMethod("get", int.class);
//                                List funcs = (List) get.invoke(dealFunc, msgCode);
//                                String logs = "deal msg code:" + msgCode + ",";
//                                logs += "func: " + funcs;
//                                log.i(logs);
//                                pkgRecvFile.write(logs);
//                            }
//                        });
//            }
//        });

//        ClassLoadReg.put(className.WriteTimeClass, new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class WriteTimeClass = XposedHelpers.findClass(clz, classLoader);
//                XposedHelpers.findAndHookMethod(WriteTimeClass, "A00", long.class, OutputStream.class,
//                        new XC_MethodHook() {
//                            @Override
//                            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                                super.beforeHookedMethod(param);
//                                log.i("write time:" + param.args[0]);
//                            }
//                        });
//            }
//        });
//        watch value
//        ClassLoadReg.put("X.0FV", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class clz0FV = XposedHelpers.findClass("X.0FV", classLoader);
//                XposedBridge.hookAllConstructors(clz0FV, new XC_MethodHook() {
//                    @Override
//                    protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                        super.afterHookedMethod(param);
//                        log.i("init msg init: " + param.thisObject);
//                        initMsg = param.thisObject;
//                        new Thread(new Runnable() {
//                            @Override
//                            public void run() {
//                                while (true){
//                                    try {
//                                        Thread.sleep(1000);
//                                    } catch (InterruptedException e) {
//                                    }
//                                    log.i("init data: " + frida_helper.object_2_string(initMsg));
//                                }
//                            }
//                        }).start();
//                    }
//                });
//            }
//        });

//        for test
//        ClassLoadReg.put("X.0D3", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class clz0D3 = XposedHelpers.findClass("X.0D3", classLoader);
////                write_list_A00
////                XposedHelpers.findAndHookMethod(clz0D3, "A00", new XC_MethodHook() {
////                    @Override
////                    protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
////                        super.beforeHookedMethod(param);
////
////                    }
////                });
////X.Type_0g8.write_any_type(X.0CM, java.lang.Object)
////                XposedHelpers.findAndHookMethod(XposedHelpers.findClass("X.0g8", classLoader), "A00",
////                        XposedHelpers.findClass("X.0CM", classLoader), Object.class,
////                        new XC_MethodHook() {
////                            @Override
////                            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
////                                super.beforeHookedMethod(param);
////                                log.i("type: " + param.args[1].getClass() + " " + param.args[1]);
////                                String[] clzName = {"X.0g9", "X.0gA", "X.0gC", "X.0gD", "X.0gE", "X.0gG", "X.0gI", "X.0gJ", "X.0gK", "X.0gL", "X.0gM", "X.0gN", "X.0gO",
////                                };
////                                Class clzClass[] = new Class[clzName.length];
////                                for (int i = 0; i < clzName.length; i++) {
////                                    clzClass[i] = XposedHelpers.findClass(clzName[i], classLoader);
////                                    log.i(clzName[i] + ":" + clzClass[i].isInstance(param.args[1]));
////                                }
////                            }
////                        });
//            }
//        });

//        ClassLoadReg.put("X.0CN", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class buff = XposedHelpers.findClass("X.0CN", classLoader);
//                XposedHelpers.findAndHookMethod(buff, "AUK", byte[].class, int.class, int.class,
//                        new XC_MethodHook() {
//                            @Override
//                            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
//                                super.beforeHookedMethod(param);
//                                new Throwable().printStackTrace();
//                            }
//                        });
//            }
//        });
//        ClassLoadReg.put("X.0fs", new ClassLoadCallBack() {
//            @Override
//            public void OnLoadedClass(String clz, ClassLoader classLoader) {
//                Class recvThread = XposedHelpers.findClass("X.0fs", classLoader);
//                XposedBridge.hookAllConstructors(recvThread, new XC_MethodHook() {
//                    @Override
//                    protected void afterHookedMethod(MethodHookParam param) throws Throwable {
//                        super.afterHookedMethod(param);
//                        Object dataStream = Reflect.GetFieldValue(recvThread, param.thisObject, "A00");
//                        log.i("dataStream is:" + dataStream.getClass() + " " + dataStream);
//                    }
//                });
//            }
//        });
