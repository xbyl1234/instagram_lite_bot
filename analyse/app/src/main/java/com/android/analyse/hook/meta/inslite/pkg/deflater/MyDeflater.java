package com.android.analyse.hook.meta.inslite.pkg.deflater;

import com.common.log;
import com.frida.frida_helper;

import java.io.ByteArrayOutputStream;
import java.util.zip.Deflater;

public class MyDeflater extends Deflater {
    ByteArrayOutputStream dataIn = new ByteArrayOutputStream();
    ByteArrayOutputStream dataOut = new ByteArrayOutputStream();


    public MyDeflater(int level) {
        super(level);
    }

    @Override
    public void setInput(byte[] b, int off, int len) {
        log.i("setInput");
        synchronized (this) {
            dataIn.write(b, off, len);
        }
        super.setInput(b, off, len);
    }

    @Override
    public int deflate(byte[] b, int off, int len, int flush) {
        log.i("deflate");
        int ret = super.deflate(b, off, len, flush);
        synchronized (this) {
            if (ret > 0) {
                dataOut.write(b, off, ret);
            }
        }
        return ret;
    }

    @Override
    public void reset() {
        log.i("reset");

        log.i("zip " + this + " in: " + frida_helper.byte_2_hex_str(dataIn.toByteArray()));
        log.i("zip " + this + " out: " + frida_helper.byte_2_hex_str(dataOut.toByteArray()));

        dataIn.reset();
        dataOut.reset();
        super.reset();
    }

}
