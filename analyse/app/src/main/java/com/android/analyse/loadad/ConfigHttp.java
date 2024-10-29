package com.android.analyse.loadad;

import com.common.log;

import org.json.JSONObject;

public class ConfigHttp {
    static HttpService server = null;
    static Response config;

    public static class Response {
        String service_ip;
        String service_port;

        public Response(String jsonStr) {
            try {
                JSONObject json = new JSONObject(jsonStr);
                if (json.has("service_port"))
                    service_port = json.getString("service_port");
                if (json.has("service_ip"))
                    service_ip = json.getString("service_ip");
            } catch (Throwable e) {
                log.e("parse Response: " + e);
            }
        }

        public String toJsonStr() {
            JSONObject json = new JSONObject();
            try {
                json.put("service_ip", service_ip);
                json.put("service_port", service_port);
            } catch (Throwable e) {
                log.e("Response: " + e);
            }
            return json.toString();
        }
    }

    public static boolean SetConfig(String json) {
        config = new Response(json);
        return true;
    }

    public static boolean start() {
        try {
            server = new HttpService("127.0.0.1", 7999);
            server.registerHandler("/config", new HttpService.HttpServerCallback() {
                @Override
                public String OnHttp(String url, JSONObject body) throws Throwable {
                    log.i("on get config: " + config.toJsonStr());
                    return config.toJsonStr();
                }
            });
            server.start();
            return true;
        } catch (Throwable e) {
            log.e("start config http error: " + e);
            return false;
        }
    }

    public static boolean stop() {
        if (server != null) {
            server.stop();
            server = null;
        }
        return true;
    }
}
