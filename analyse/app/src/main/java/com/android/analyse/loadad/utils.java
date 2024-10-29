package com.android.analyse.loadad;

import com.common.log;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLEncoder;
import java.util.Map;

public class utils {

    public static String httpGet(String urlString, Map<String, String> params) {
        return httpGet(urlString, params, 5);
    }

    public static String httpGet(String urlString, Map<String, String> params, int timeout) {
        try {
            if (params != null && !params.isEmpty()) {
                log.i("httpGet1: " + urlString);
                int index = 0;
                for (Map.Entry<String, String> entry : params.entrySet()) {

                    log.i("httpGet2: " + entry.getKey() + "---" + entry.getValue());
                    if (index == 0) {
                        urlString = urlString + "?" + entry.getKey() + "=" + URLEncoder.encode(entry.getValue(), "UTF-8");
                    } else {
                        urlString = urlString + "&" + entry.getKey() + "=" + URLEncoder.encode(entry.getValue(), "UTF-8");
                    }
                    index++;
                }
            }

            log.i("httpGet: " + urlString);
            URL url = new URL(urlString);
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("GET");
            connection.setConnectTimeout(timeout * 1000);
            connection.setReadTimeout(timeout * 1000);
            log.i("http response code: " + connection.getResponseCode());
            BufferedReader reader = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            StringBuilder response = new StringBuilder();
            String line;
            while ((line = reader.readLine()) != null) {
                response.append(line);
            }
            reader.close();
            return response.toString();
        } catch (Exception e) {
            log.i("http request error: " + e);
            e.printStackTrace();
            return null;
        }
    }
}
