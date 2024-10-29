package com.android.analyse.hook.meta.inslite.pkg.deflater;

import java.io.IOException;
import java.io.OutputStream;
import java.util.zip.Deflater;
import java.util.zip.DeflaterOutputStream;

public class MyDeflaterOutputStream extends DeflaterOutputStream {
    public MyDeflaterOutputStream(OutputStream out, Deflater def, int size) throws IOException {
        super(out, def, size);
    }

    @Override
    public void write(byte[] b, int off, int len) throws IOException {
        super.write(b, off, len);
    }
}
