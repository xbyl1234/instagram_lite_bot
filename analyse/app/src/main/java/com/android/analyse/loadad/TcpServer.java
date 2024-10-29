package com.android.analyse.loadad;

import android.util.Log;


import java.io.InputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;

public class TcpServer extends TcpChannel {
    private volatile boolean started = false;
    private volatile AtomicLong id = new AtomicLong(0);
    private final List<Socket> clientSockets = new ArrayList<>();
    private final Map<Long, Consumer<String>> callbacks = new ConcurrentHashMap<>();

    private static final TcpServer INSTANCE = new TcpServer();

    private volatile ServerSocket serverSocket = null;

    private TcpServer() {
        super();
    }

    public static TcpServer getInstance() {
        return INSTANCE;
    }

    public void start() {
        if (started) return;

        started = true;
        stopped = false;
        clientSockets.clear();

        executorService.submit(this::acceptConnections);
    }

    public void stop() {
        started = false;
        stopped = true;

        synchronized (clientSockets) {
            clientSockets.forEach(this::closeSocket);
            clientSockets.clear();
        }
        if (serverSocket != null) {
            try {
                serverSocket.close();
                serverSocket = null;
            } catch (Exception e) {
            }
        }
    }

    public void send(int requestCode, String params, Consumer<String> callback) {
        if (clientSockets.isEmpty()) {
            Log.d("Adapter", "没有任何客户端存在");
            callback.accept("error");
            return;
        }
        long idVal = id.incrementAndGet();

        callbacks.put(idVal, callback);

        String message = Util.toJson(TcpMessage.of(requestCode, params, idVal));
        byte[] messageBytes = message.getBytes();

        short length = (short) messageBytes.length;
        byte[] lengthBytes = ByteBuffer.allocate(2).putShort(length).array();

        synchronized (clientSockets) {
            clientSockets.forEach(socket -> sendPackage(socket, lengthBytes, messageBytes));
        }
        new Thread(() -> {
            try {
                Thread.sleep(5000);
                if (!stopped) {
                    synchronized (clientSockets) {
                        Consumer<String> back = callbacks.getOrDefault(idVal, null);
                        if (back != null) {
                            back.accept("error");
                            callbacks.remove(idVal);
                        }
                    }
                }
            } catch (Exception e) {
            }
        }).start();
    }

    private void acceptConnections() {
        Log.d("Adapter", "Start to accept connections...");
        try {
            serverSocket = new ServerSocket(9003);
            while (!stopped) {
                Socket clientSocket = serverSocket.accept();

                Log.d("Adapter", "Accept connection from " + clientSocket.getInetAddress().toString() + " " + clientSocket.getPort());

                executorService.submit(() -> processClientMessages(clientSocket));
            }
        } catch (Exception exception) {
            exception.printStackTrace();
        }
    }

    private void processClientMessages(Socket clientSocket) {
        try {
            Log.d("Adapter", "Start receive message from " + clientSocket.getInetAddress().toString());

            synchronized (clientSockets) {
                clientSockets.add(clientSocket);
            }

            // 获取输入流和输出流
            InputStream inputStream = clientSocket.getInputStream();

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
                    TcpMessage response = Util.fromJson(message, TcpMessage.class);
                    Consumer<String> callback = callbacks.getOrDefault(response.id, null);
                    if (callback != null) {
                        callback.accept(response.data);
                    }
                    callbacks.remove(response.id);
                }
            } while (!stopped);

            Log.d("Adapter", "Server side finished");
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            // 关闭连接
            closeSocket(clientSocket);
            System.out.println("客户端已断开连接：" + clientSocket.getInetAddress().getHostAddress());

            synchronized (clientSockets) {
                clientSockets.remove(clientSocket);
            }
        }
    }
}