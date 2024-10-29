package com.android.analyse.loadad;

import android.util.Log;


import java.io.IOException;
import java.io.InputStream;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;
import java.nio.ByteBuffer;

public class TcpClient extends TcpChannel {
    public interface IMessageCallback {
        void onMessage(int requestCode, String data);
    }

    private volatile boolean started = false;
    private volatile IMessageCallback callback = null;
    private volatile Socket socket;

    private static final TcpClient INSTANCE = new TcpClient();

    private TcpClient() {
        super();
    }

    public static TcpClient getInstance() {
        return INSTANCE;
    }

    public void start(IMessageCallback callback) {
        if (started) return;

        started = true;
        stopped = false;

        if (socket != null) {
            closeSocket(socket);
            socket = null;
        }

        connect("localhost", 9003, callback);
    }

    public void stop() {
        started = false;
        stopped = true;

        closeSocket(socket);
        socket = null;
    }

    private boolean connect(String ip, int port, IMessageCallback callback) {
        Log.d("Adapter", "Start connecting to server...");

        socket = new Socket();
        try {

            InetAddress inetAddress = InetAddress.getByName(ip);
            SocketAddress socketAddress = new InetSocketAddress(inetAddress, port);
            socket.connect(socketAddress, 10000);

            if (socket.isConnected()) {
                Log.d("Adapter", "Connected to server...");

                this.callback = callback;

                started = true;
                stopped = false;

                executorService.submit(new Runnable() {
                    @Override
                    public void run() {
                        processServerMessages();
                    }
                });

                return true;
            }

            Log.d("Adapter", "Connect to server failed.");

            started = false;
            return false;
        } catch (IOException e) {
            Log.d("Adapter", "Connect to server failed.");

            e.printStackTrace();
            return false;
        }
    }

    public boolean send(long id, int requestCode, String params) {
        if (socket == null || !socket.isConnected()) {
            Log.d("Adapter", "客户端连接服务器断开");
            return false;
        }

        String message = TcpMessage.of(requestCode, params, id).toJsonStr();
        byte[] messageBytes = message.getBytes();

        short length = (short) messageBytes.length;
        byte[] lengthBytes = ByteBuffer.allocate(2).putShort(length).array();

        sendPackage(socket, lengthBytes, messageBytes);

        return true;
    }

    private void processServerMessages() {
        try {
            // 获取输入流和输出流
            InputStream inputStream = socket.getInputStream();

            // 读取客户端发送的数据
            byte[] headerBuffer = new byte[2];
            byte[] messageBuffer = new byte[65536];

            do {
                // 读头
                int packetLength = readHeader(inputStream, headerBuffer);
                if (packetLength < 0) {
                    break;
                }

                String message = readMessage(inputStream, messageBuffer, packetLength);

                if (message != null) {
                    TcpMessage response = new TcpMessage(message);
                    if (callback != null) {
                        callback.onMessage(response.requestCode, response.data);
                    }
                }
            } while (!stopped);

            Log.d("Adapter", "Client side finished");
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            // 关闭连接
            closeSocket(socket);
            System.out.println("和服务器之间的连接已断开：" + socket.getInetAddress().getHostAddress());

            socket = null;
            stopped = true;
            started = false;
        }
    }
}
