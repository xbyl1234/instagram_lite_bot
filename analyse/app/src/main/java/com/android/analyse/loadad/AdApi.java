package com.android.analyse.loadad;

import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;

import com.common.log;

import org.json.JSONObject;

import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Consumer;

public class AdApi {
    private static final String gameUrl = "http://localhost:9001/";


    static public class GeneralData {
        public int code;
        public String data;
        public String msg;

        public GeneralData(String jsonStr) {
            try {
                JSONObject json = new JSONObject(jsonStr);
                if (json.has("code"))
                    code = json.getInt("code");
                if (json.has("data"))
                    data = json.getString("data");
                if (json.has("msg"))
                    msg = json.getString("msg");
            } catch (Throwable e) {
                log.e("parse GeneralData: " + e);
            }
        }

        public String toJsonString() {
            JSONObject json = new JSONObject();
            try {
                json.put("code", code);
                json.put("data", data);
                json.put("msg", msg);
            } catch (Throwable e) {
                log.e("GeneralData2str: " + e);
            }
            return json.toString();
        }
    }

    public static void openApp(Context context, String packageName, String activity) {
        ComponentName component = new ComponentName(packageName, activity);
        Intent intent = new Intent();
        intent.setComponent(component);
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        intent.putExtra("url", "127.0.0.1");
        intent.putExtra("port", 9003);
        context.startActivity(intent);
    }

    private static boolean isResponseOk(String json) {
        try {
            if (json != null && json != "") {
                GeneralData data = new GeneralData(json);
                return data != null && data.code == 0;
            } else {
                return false;
            }
        } catch (Exception e) {
            return false;
        }
    }

    public static boolean playVideo(boolean useHttp) {
        if (useHttp) {
            String resp = utils.httpGet(gameUrl + "playvideoad", null);
            log.i("play video http response: " + resp);
            return isResponseOk(resp);
        } else {
            AtomicReference<String> result = new AtomicReference<>();
            TcpServer.getInstance().send(1003, null, new Consumer<String>() {
                @Override
                public void accept(String resp) {
                    result.set(resp);
                }
            });
            while (true) {
                if (result.get() != null && result.get() != "") {
                    break;
                }
                try {
                    Thread.sleep(200);
                } catch (InterruptedException e) {
                }
            }
            log.i("play video socket response: " + result.get());
            return "0".equals(result.get());
        }
    }
}
