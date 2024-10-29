package com.android.analyse.loadad;

import android.util.Log;

import com.common.log;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class TcpChannel {
    protected volatile boolean stopped = false;
    protected final ExecutorService executorService;

    protected TcpChannel() {
        executorService = Executors.newFixedThreadPool(10);
    }

    protected void closeSocket(Socket socket) {
        try {
            if (socket != null) {
                socket.close();
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    protected void sendPackage(Socket socket, byte[] lengthBytes, byte[] messageBytes) {
        try {
            OutputStream outputStream = socket.getOutputStream();
            outputStream.write(lengthBytes);
            outputStream.write(messageBytes);
            outputStream.flush();
        } catch (Exception exception) {
            exception.printStackTrace();
        }
    }

    protected int readHeader(InputStream inputStream, byte[] buffer) throws IOException {
        int bytesRead = 0;
        int headerBufferRead = 0;
        int headerSize = 2;

        while (!stopped && (bytesRead = inputStream.read(buffer, headerBufferRead, headerSize - headerBufferRead)) != -1) {
            headerBufferRead += bytesRead;

            if (headerBufferRead >= headerSize) {
                ByteBuffer byteBuffer = ByteBuffer.allocate(Integer.BYTES);
                byteBuffer.put(buffer);
                byteBuffer.rewind();

                short result = byteBuffer.getShort();

                Log.d("Adapter", "Read package length: " + result);
                return result;
            }
        }

        return -1;
    }

    protected String readMessage(InputStream inputStream, byte[] buffer, int packageLength) throws IOException {
        if (packageLength <= 0) {
            return null;
        }

        int bytesRead = 0;
        int messageBufferRead = 0;

        while (!stopped && (bytesRead = inputStream.read(buffer, messageBufferRead, packageLength - messageBufferRead)) != -1) {
            messageBufferRead += bytesRead;

            if (messageBufferRead >= packageLength) {
                return new String(buffer, 0, packageLength, StandardCharsets.UTF_8);
            }
        }

        return null;
    }
}
