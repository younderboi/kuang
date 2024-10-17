package agent

import (
	"bufio"
	"fmt"
	"net"
)

type AgentTransport interface {
	Connect() error
	Read() (string, error)
	Write(data string) error
	Close() error
}

// === TCP Transport
type TCPTransport struct {
	LHOST string
	LPORT string
	conn  net.Conn
}

func (t *TCPTransport) Connect() error {
	var err error
	t.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", t.LHOST, t.LPORT))
	return err
}

func (t *TCPTransport) Read() (string, error) {
	reader := bufio.NewReader(t.conn)
	return reader.ReadString('\n')
}

func (t *TCPTransport) Write(data string) error {
	_, err := t.conn.Write([]byte(data + "\n"))
	return err
}

func (t *TCPTransport) Close() error {
	return t.conn.Close()
}
