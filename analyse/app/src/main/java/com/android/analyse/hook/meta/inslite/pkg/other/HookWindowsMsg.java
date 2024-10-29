package com.android.analyse.hook.meta.inslite.pkg.other;

import android.os.Message;

import com.android.analyse.hook.AppFileWriter;
import com.android.analyse.hook.meta.common.ClassLoadCallBack;
import com.android.analyse.hook.meta.inslite.FieldName;
import com.common.tools.hooker.HookTools;

import java.util.HashMap;
import java.util.Map;

import de.robv.android.xposed.XC_MethodHook;
import de.robv.android.xposed.XposedHelpers;

public class HookWindowsMsg extends ClassLoadCallBack {
    public HookWindowsMsg(AppFileWriter logFile) {
        super(logFile);
    }

    static class MsgInfo {
        int index;
        String name;

        public MsgInfo(int index, String name) {
            this.index = index;
            this.name = name;
        }
    }

    static Map<Integer, MsgInfo> msgMaps = new HashMap<Integer, MsgInfo>() {{
        put(0, new MsgInfo(0, "NO_TYPE"));
        put(1, new MsgInfo(1, "CONNECTING"));
        put(2, new MsgInfo(2, "OBSOLETE_CONNECTIONDISCONNECTED"));
        put(3, new MsgInfo(3, "CONNECTIONESTABLISHED"));
        put(4, new MsgInfo(4, "MESSAGERECEIVED"));
        put(7, new MsgInfo(5, "CONNECTIONNOTAUTHORIZED"));
        put(8, new MsgInfo(6, "CONNECTORFAILURE"));
        put(10, new MsgInfo(7, "KEYPRESSED"));
        put(11, new MsgInfo(8, "KEYREPEATED"));
        put(12, new MsgInfo(9, "POINTERPRESSED"));
        put(13, new MsgInfo(10, "POINTERRELEASED"));
        put(14, new MsgInfo(11, "POINTERDRAGGED"));
        put(15, new MsgInfo(12, "NONPRIMARYPTR_PRESSED"));
        put(20, new MsgInfo(13, "REPORTEXCEPTION"));
        put(21, new MsgInfo(14, "REPORTLOG"));
        put(22, new MsgInfo(15, "REPORTLOG2"));
        put(23, new MsgInfo(16, "REPORTLOG3"));
        put(31, new MsgInfo(17, "TIMEOUT"));
        put(32, new MsgInfo(18, "OBSOLETE_CONNECTION_ATTEMPT"));
        put(40, new MsgInfo(19, "RUNNABLE"));
        put(41, new MsgInfo(20, "OBSOLETE_FORCERELOGIN"));
        put(60, new MsgInfo(21, "UI_RESOLUTION_CHANGED"));
        put(61, new MsgInfo(22, "OBSOLETE_ENTER_STANDBY_MODE"));
        put(62, new MsgInfo(23, "PERFORM_SCREEN_ACTION"));
        put(63, new MsgInfo(24, "PUSH_TOKEN_RECEIVED"));
        put(64, new MsgInfo(25, "PUSH_ERROR_RECEIVED"));
        put(65, new MsgInfo(26, "PUSH_PAYLOAD_RECEIVED"));
        put(67, new MsgInfo(27, "APP_REQUEST_RECEIVED"));
        put(68, new MsgInfo(28, "OBSOLETE_PUSH_STATUS_RECEIVED"));
        put(70, new MsgInfo(29, "LOAD_IMAGE_RESOURCE"));
        put(71, new MsgInfo(30, "HTTP_REQUEST_REPLY"));
        put(72, new MsgInfo(31, "OBSOLETE_SEND_INSTRUMENT_DATA"));
        put(74, new MsgInfo(32, "OBSOLETE_RECORD_INSTRUMENT_VALUE"));
        put(76, new MsgInfo(33, "SCROLL_UPDATE"));
        put(77, new MsgInfo(34, "SECOND_CLICK_TIMEOUT"));
        put(78, new MsgInfo(35, "LONG_CLICK_TIMEOUT"));
        put(79, new MsgInfo(36, "POINTERRELEASED_SUPPRESS_ACTION"));
        put(81, new MsgInfo(37, "SMARTPHONE_PUSH_TOKEN_RECEIVED"));
        put(83, new MsgInfo(38, "DECODED_IMAGE_FROM_CACHE"));
        put(85, new MsgInfo(39, "NETWORK_CONNECTED"));
        put(86, new MsgInfo(40, "APP_DESTROY"));
        put(87, new MsgInfo(41, "NETWORK_DISCONNECTED"));
        put(89, new MsgInfo(42, "OBSOLETE_FIRST_SCREEN_RECEIVED"));
        put(91, new MsgInfo(43, "IMAGE_CALC_INSTRUMENT"));
        put(92, new MsgInfo(44, "SIMPLE_TOUCH"));
        put(207, new MsgInfo(45, "LOCATION_FETCHED"));
        put(209, new MsgInfo(46, "OBSOLETE_WAKE_FROM_STANDBY"));
        put(210, new MsgInfo(47, "ANIMATE_STARTUP_SCREEN"));
        put(211, new MsgInfo(48, "OBSOLETE_SEND_CACHE_PERFORMANCE"));
        put(212, new MsgInfo(49, "VIDEO_DOWNLOAD_COMPLETE"));
        put(213, new MsgInfo(50, "SMS_CONFIRMATION_CODE_RECEIVED"));
        put(214, new MsgInfo(51, "OAUTH_TOKEN_RECEIVED"));
        put(215, new MsgInfo(52, "DATA_RECEIVED"));
        put(216, new MsgInfo(53, "DATE_PICKED"));
        put(218, new MsgInfo(54, "RESET_SESSION"));
        put(219, new MsgInfo(55, "SNAPPED_IN_HSCROLL"));
        put(221, new MsgInfo(56, "RELOGIN_TO_SERVER"));
        put(222, new MsgInfo(57, "UPLOAD_VIDEO_TO_THREAD"));
        put(226, new MsgInfo(58, "OAUTH_TOKENS_RECEIVED"));
        put(227, new MsgInfo(59, "WIDGET_PAYLOAD_RECEIVED"));
        put(228, new MsgInfo(60, "SHOW_CLEAR_WWW_ROUTING_BUTTON"));
    }};

    static public String MsgObj2Str(Object obj) {
        String logs = "";
        logs += "msg: " + obj + ", ";
        logs += "type: " + HookTools.GetFieldValue(obj.getClass(), obj, FieldName.HookWindowsMsg_type) + ", ";
        logs += "sub_type: " + HookTools.GetFieldValue(obj.getClass(), obj, FieldName.HookWindowsMsg_sub_type) + ", ";
        Object screen = HookTools.GetFieldValue(obj.getClass(), obj, FieldName.HookWindowsMsg_screen);
        if (screen != null) {
            logs += "screen_id: " + HookTools.GetFieldValue(screen.getClass(), screen, FieldName.HookWindowsMsg_screen_id) + ", ";
        }
        return logs;
    }

    @Override
    public void OnLoadedClass(String clzName, Class clz, ClassLoader classLoader) {
        XposedHelpers.findAndHookMethod(clz, "handleMessage", Message.class, new XC_MethodHook() {
            @Override
            protected void beforeHookedMethod(MethodHookParam param) throws Throwable {
                super.beforeHookedMethod(param);
                String logs = "handleMessage ";
                Message msg = (Message) param.args[0];
                logs += MsgObj2Str(msg.obj);
                write_log(logs);
            }
        });
    }
}
