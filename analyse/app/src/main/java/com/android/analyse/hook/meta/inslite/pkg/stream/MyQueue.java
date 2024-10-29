package com.android.analyse.hook.meta.inslite.pkg.stream;

import com.android.analyse.hook.meta.inslite.pkg.msg.Message;
import com.common.log;

import java.util.concurrent.LinkedBlockingQueue;

public class MyQueue extends LinkedBlockingQueue {
    public MyQueue() {
    }

    @Override
    public boolean offer(Object e) {
        log.i("MyQueue put: " + Message.Msg2Json(e,true));
        return super.offer(e);
    }
}
