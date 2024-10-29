package com.android.analyse.hook.meta.inslite.pkg.msg;

import com.alibaba.fastjson2.JSONObject;
import com.android.analyse.hook.meta.inslite.ClassNames;
import com.android.analyse.hook.meta.inslite.FieldName;
import com.android.analyse.hook.meta.common.ClassloaderHook;
import com.common.tools.hooker.HookTools;
import com.frida.frida_helper;

public class Message {
    static void SlicedByteBuffer2Json(JSONObject json, Object msg, boolean needData) {
        int len = (int) HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.slicedByteBuffer), msg, FieldName.SlicedByteBuffer_AllSizeName);
        json.put("all_size", len);
        json.put("cur_size", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.slicedByteBuffer), msg, FieldName.SlicedByteBuffer_CurSizeName));
        if (needData) {
            json.put("data", frida_helper.byte_2_hex_str(HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.slicedByteBuffer), msg, FieldName.SlicedByteBuffer_BuffName), 0, len));
        }
    }

    static void SendMsgData2Json(JSONObject json, Object msg) {
        json.put("clz_name", msg.getClass().getName());
        json.put("msg_code", Message.GetMsgCode(msg));
        json.put("sender_index", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgData), msg, FieldName.HookSendMsgData_sender_index));
        json.put("recver_index", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgData), msg, FieldName.HookSendMsgData_recver_index));
        json.put("is_zip", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgData), msg, FieldName.HookSendMsgData_is_zip));
        json.put("unknow_bool", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgData), msg, FieldName.HookSendMsgData_unknow_bool));
        json.put("unknow_long", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgData), msg, FieldName.HookSendMsgData_unknow_long));
    }

    static void SendMsgChildData2Json(JSONObject json, Object obj) {
        json.put("child_int", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgDataCld), obj, FieldName.SendMsgChild_child_int));
        json.put("child_long", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgDataCld), obj, FieldName.SendMsgChild_child_long));
    }

    public static void RecvMsgData2Json(JSONObject json, Object msg) {
        json.put("msg_code", Message.GetMsgCode(msg));
    }

    public static int GetMsgCode(Object msg) {
        try {
            Class msgClz = msg.getClass();
            if (msgClz == ClassloaderHook.GetClass(ClassNames.sendMsgDataCld) ||
                    msgClz == ClassloaderHook.GetClass(ClassNames.sendMsgData)) {
                return (int) HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.sendMsgData), msg, FieldName.HookSendMsgData_msg_code);
            } else if (msgClz == ClassloaderHook.GetClass(ClassNames.recvMsg)) {
                return (int) HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.recvMsg), msg, FieldName.RecvMsg_MsgCodeName);
            }
        } catch (Throwable e) {
            return 0;
        }
        return 0;
    }

    public static JSONObject MsgOffset(Object msg) {
        JSONObject json = new JSONObject();
        int len = (int) HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.slicedByteBuffer), msg, FieldName.SlicedByteBuffer_AllSizeName);
        json.put("all_size", len);
        json.put("cur_size", HookTools.GetFieldValue(ClassloaderHook.GetClass(ClassNames.slicedByteBuffer), msg, FieldName.SlicedByteBuffer_CurSizeName));
        return json;
    }

    public static JSONObject Msg2Json(Object msg, boolean needData) {
        if (msg == null) {
            return new JSONObject();
        }
        JSONObject json = new JSONObject();
        Class msgClz = msg.getClass();
        if (msgClz == ClassloaderHook.GetClass(ClassNames.sendMsgDataCld)) {
            SendMsgChildData2Json(json, msg);
            SendMsgData2Json(json, msg);
            SlicedByteBuffer2Json(json, msg, needData);
        } else if (msgClz == ClassloaderHook.GetClass(ClassNames.sendMsgData)) {
            SendMsgData2Json(json, msg);
            SlicedByteBuffer2Json(json, msg, needData);
        } else if (msgClz == ClassloaderHook.GetClass(ClassNames.slicedByteBuffer)) {
            SlicedByteBuffer2Json(json, msg, needData);
        } else if (msgClz == ClassloaderHook.GetClass(ClassNames.recvMsg)) {
            RecvMsgData2Json(json, msg);
            SlicedByteBuffer2Json(json, msg, needData);
        }
        return json;
    }
}