package bufferred_conn

import (
	"github.com/mzz2017/softwind/pkg/zeroalloc/bufio"
	"net"
)

type BufferedConn struct {
	r        *bufio.Reader
	net.Conn // So that most methods are embedded
}

func NewBufferedConn(c *net.TCPConn) BufferedConn {
	return BufferedConn{bufio.NewReader(c), c}
}

func NewBufferedConnSize(c *net.TCPConn, n int) BufferedConn {
	return BufferedConn{bufio.NewReaderSize(c, n), c}
}

func (b BufferedConn) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b BufferedConn) Close() error {
	b.r.Put()
	return b.Conn.Close()
}

func (b BufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}
