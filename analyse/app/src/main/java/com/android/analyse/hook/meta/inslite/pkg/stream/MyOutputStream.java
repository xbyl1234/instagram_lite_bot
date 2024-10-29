package com.android.analyse.hook.meta.inslite.pkg.stream;

import com.android.analyse.hook.AppFileWriter;
import com.frida.frida_helper;

import java.io.IOException;
import java.io.OutputStream;

public class MyOutputStream extends OutputStream {
    OutputStream stream;
    AppFileWriter writer;

    public MyOutputStream(OutputStream stream, AppFileWriter writer) {
        this.stream = stream;
        this.writer = writer;
    }

    @Override
    public void write(int b) throws IOException {
        byte[] fuck = new byte[1];
        fuck[0] = (byte) (b & 0xff);
        write(fuck, 0, 1);
    }

    @Override
    public void write(byte[] b) throws IOException {
        write(b, 0, b.length);
    }

    @Override
    public void write(byte[] b, int off, int len) throws IOException {
        String data = "output stream " + this + " write: " + frida_helper.byte_2_hex_str(b, off, len);
//        log.i(data);
        writer.write(data);
        stream.write(b, off, len);
    }
}
