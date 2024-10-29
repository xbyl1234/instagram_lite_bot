package net

import (
	"CentralizedControl/common"
	"CentralizedControl/common/proxys"
	"fmt"
	tls "github.com/refraction-networking/utls"
	"net"
	"time"
)

func CreateGoTls(host, port string, p proxys.Proxy) (net.Conn, error) {
	config := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true,
	}
	var conn net.Conn
	var err error
	if p != nil {
		conn, err = p.GetDialer().Dial("tcp", host+":"+port)
	} else {
		conn, err = net.Dial("tcp", host+":"+port)
	}
	if err != nil {
		return nil, err
	}
	return tls.Client(conn, config), err
}

func forgeConn(host, port string) {
	clientTcp, err := net.DialTimeout("tcp", host+":"+port, 10*time.Second)
	if err != nil {
		fmt.Printf("net.DialTimeout error: %+v", err)
		return
	}
	clientUtls := tls.UClient(clientTcp, nil, tls.HelloGolang)
	defer clientUtls.Close()
	clientUtls.SetSNI(host)
	err = clientUtls.Handshake()
	if err != nil {
		fmt.Printf("clientUtls.Handshake() error: %+v", err)
	}
	serverConn, clientConn := net.Pipe()
	clientUtls.SetUnderlyingConn(clientConn)
	hs := clientUtls.HandshakeState
	serverTls := tls.MakeConnWithCompleteHandshake(serverConn, hs.ServerHello.Vers, hs.ServerHello.CipherSuite,
		hs.MasterSecret, hs.Hello.Random, hs.ServerHello.Random, false)
	if serverTls == nil {
		fmt.Printf("tls.MakeConnWithCompleteHandshake error, unsupported TLS protocol?")
		return
	}
	go func() {
		clientUtls.Write([]byte("Hello, world!"))
		resp := make([]byte, 13)
		read, err := clientUtls.Read(resp)
		if err != nil {
			fmt.Printf("error reading client: %+v\n", err)
		}
		fmt.Printf("Client read %d bytes: %s\n", read, string(resp))
		fmt.Println("Client closing...")
		clientUtls.Close()
		fmt.Println("client closed")
	}()
	buf := make([]byte, 13)
	read, err := serverTls.Read(buf)
	if err != nil {
		fmt.Printf("error reading server: %+v\n", err)
	}
	fmt.Printf("Server read %d bytes: %s\n", read, string(buf))
	serverTls.Write([]byte("Test response"))
	serverTls.Read(buf)
	fmt.Println("Server closed")
}

func CreateJa3Socket(host, port string, p proxys.Proxy) (net.Conn, error) {
	//forgeConn(host, port)
	config := tls.Config{ServerName: host}
	var conn net.Conn
	var err error
	if p != nil {
		conn, err = p.GetDialer().Dial("tcp", host+":"+port)
	} else {
		conn, err = net.Dial("tcp", host+":"+port)
	}
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout error: %+v", err)
	}
	uTlsConn := tls.UClient(conn, &config, tls.HelloCustom)
	err = uTlsConn.ApplyPreset(
		common.GenTlsConfig("771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-49172-156-157-47-53,0-23-65281-10-11-35-5-13-51-45-43-21,29-23-24,0",
			common.GetInsLiteExtensionConfig()))
	if err != nil {
		return uTlsConn, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}
	err = uTlsConn.Handshake()
	if err != nil {
		return uTlsConn, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}
	return uTlsConn, nil
}
