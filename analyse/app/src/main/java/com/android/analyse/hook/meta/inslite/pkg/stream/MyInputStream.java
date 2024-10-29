package com.android.analyse.hook.meta.inslite.pkg.stream;

import com.android.analyse.hook.AppFileWriter;
import com.common.log;
import com.frida.frida_helper;

import java.io.IOException;
import java.io.InputStream;

public class MyInputStream extends InputStream {
    InputStream stream = null;
    String tags = null;
    AppFileWriter writer = null;

    public MyInputStream(InputStream stream, AppFileWriter writer, String tags) {
        this.stream = stream;
        this.tags = tags;
        this.writer = writer;
    }

    @Override
    public int read() throws IOException {
        byte[] result = new byte[1];
        if (read(result, 0, 1) != 1) {
            log.e("read fuck error!");
        }
        return result[0] & 0xff;
    }

    @Override
    public int read(byte[] b, int off, int len) throws IOException {
        int ret = stream.read(b, off, len);
        String logs = this.tags + " input stream " + this + " read: " + frida_helper.byte_2_hex_str(b, off, ret);
//        log.i(logs);
        writer.write(logs);
        return ret;
    }

}
