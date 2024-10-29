package com.android.analyse;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import android.view.View;

import com.android.analyse.databinding.ActivityMainBinding;
import com.android.analyse.loadad.AdApi;
import com.android.analyse.loadad.TcpServer;

public class MainActivity extends AppCompatActivity {
    private ActivityMainBinding binding;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        binding = ActivityMainBinding.inflate(getLayoutInflater());
        setContentView(binding.getRoot());
        FileManger.CpyFile(this);
        findViewById(R.id.button).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
//                System.loadLibrary("l2f1abd97");
//                new Thread(new Runnable() {
//                    @Override
//                    public void run() {
//                        TcpServer.getInstance().start();
//                        AdApi.openApp(MainActivity.this, "com.wallpapers.explorers.lop", "a.C");
//                        AdApi.openApp(MainActivity.this, "com.myphone.myscreen", "fa.df.af.CCSAG");
//                        while (true) {
//                            AdApi.playVideo(false);
//                            try {
//                                Thread.sleep(1500);
//                            } catch (Throwable e) {
//                            }
//                        }
//                    }
//                }).start();
            }
        });


//        Test.Test();
//        System.loadLibrary("analyse");
//        Native.nativeInitJniTrace(frida_helper.class);
//        String dataPaht = "/data/data/" + getApplicationInfo().packageName;
//        nativeDumpSo("libc.so", dataPaht);
//        nativeInitHook();
//        nativeEnumSymbols();
//        System.loadLibrary("test");
    }
}