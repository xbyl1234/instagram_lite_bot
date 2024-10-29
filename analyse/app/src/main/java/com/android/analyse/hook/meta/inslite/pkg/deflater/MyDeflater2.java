package com.android.analyse.hook.meta.inslite.pkg.deflater;

import com.android.analyse.hook.AppFileWriter;
import com.common.log;
import com.frida.frida_helper;

import java.util.zip.Deflater;

public class MyDeflater2 extends Deflater {
    AppFileWriter file;

    public MyDeflater2(int level, AppFileWriter file) {
        super(level);
        this.file = file;
    }

    @Override
    public void setInput(byte[] b, int off, int len) {
        log.i("zip " + this + " setInput: " + frida_helper.byte_2_hex_str(b, off, len));
        super.setInput(b, off, len);
    }

    @Override
    public int deflate(byte[] b, int off, int len, int flush) {
        int ret = super.deflate(b, off, len, flush);
        if (ret > 0) {
            String logs = "zip " + this + " deflate: " + frida_helper.byte_2_hex_str(b, off, ret);
            log.i(logs);
            file.write(logs);
        }
        return ret;
    }

}
