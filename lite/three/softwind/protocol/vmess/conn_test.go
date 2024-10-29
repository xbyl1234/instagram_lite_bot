package vmess

import (
	"bufio"
	"crypto/tls"
	"github.com/google/uuid"
	"github.com/mzz2017/softwind/protocol"
	"io"
	"net"
	"net/http"
	"testing"
)

func TestClientConn(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:18080")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	id, err := uuid.Parse("28446de9-2a7e-4fab-827b-6df93e46f945")
	if err != nil {
		t.Fatal(err)
	}
	conn, err = NewConn(conn, Metadata{
		Metadata:protocol.Metadata{
			Type:     protocol.MetadataTypeDomain,
			Hostname: "www.qq.com",
			Port:     443,
			Cipher:   string(CipherC20P1305),
			IsClient: true,
		},
		Network:  "tcp",
	}, NewID(id).CmdKey())
	defer conn.Close()
	conn = tls.Client(conn, &tls.Config{
		ServerName: "www.qq.com",
	})
	defer conn.Close()
	req, _ := http.NewRequest("GET", "https://www.qq.com", nil)
	if err := req.Write(conn); err != nil {
		t.Fatal(err)
	}
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Header, resp.Status, string(b))
	conn.Write(nil)
}
