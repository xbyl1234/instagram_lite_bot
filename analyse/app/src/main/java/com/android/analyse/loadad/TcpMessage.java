package com.android.analyse.loadad;

import org.json.JSONObject;

public class TcpMessage {
    public int requestCode;
    public String data;
    public long id = 0;

    public static TcpMessage of(int requestCode, String data, long id) {
        TcpMessage tcpMessage = new TcpMessage();
        tcpMessage.requestCode = requestCode;
        tcpMessage.data = data;
        tcpMessage.id = id;
        return tcpMessage;
    }

    public TcpMessage() {
    }

    public TcpMessage(String jsonStr) {
        try {
            JSONObject json = new JSONObject(jsonStr);
            requestCode = json.getInt("requestCode");
            data = json.getString("data");
            id = json.getLong("id");
        } catch (Throwable e) {
            e.printStackTrace();
        }
    }

    public String toJsonStr() {
        JSONObject json = new JSONObject();
        try {
            json.put("requestCode", requestCode);
            json.put("data", data);
            json.put("id", id);
        } catch (Throwable e) {
            e.printStackTrace();
        }
        return json.toString();
    }


}
