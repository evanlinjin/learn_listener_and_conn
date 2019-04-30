package netutil

import (
	"errors"
	"io"
	"net"
	"net/rpc"
)

const CacheCount = 100
const BufferSize = 1024

type PrefixConn struct {
	net.Conn             // original conn.
	inCh     chan []byte // original conn pushes here.
	prefix   byte        // Prefix to append before sending to original conn.
}

func (s *PrefixConn) Read(p []byte) (n int, err error) {
	data := <-s.inCh
	if len(data) > len(p) {
		return 0, io.ErrShortBuffer
	}
	copy(p, data)
	return len(data), nil
}

func (s *PrefixConn) Write(p []byte) (n int, err error) {
	data := make([]byte, len(p)+1)
	data[0] = s.prefix
	copy(data[1:], p)
	n, err = s.Conn.Write(data)
	return n - 1, err
}

type RPCDuplex struct {
	origin net.Conn
	server *PrefixConn
	client *PrefixConn
	rpcS   *rpc.Server
	rpcC   *rpc.Client
}

func NewRPCDuplex(conn net.Conn, srv *rpc.Server, init bool) *RPCDuplex {
	var serverPrefix, clientPrefix byte
	if init {
		serverPrefix, clientPrefix = 0, 1
	} else {
		serverPrefix, clientPrefix = 1, 0
	}
	serverConn := &PrefixConn{
		Conn:   conn,
		inCh:   make(chan []byte, CacheCount),
		prefix: serverPrefix,
	}
	clientConn := &PrefixConn{
		Conn:   conn,
		inCh:   make(chan []byte, CacheCount),
		prefix: clientPrefix,
	}
	d := &RPCDuplex{
		origin: conn,
		server: serverConn,
		client: clientConn,
		rpcS:   srv,
		rpcC:   rpc.NewClient(clientConn),
	}
	return d
}

func (d *RPCDuplex) Serve() error {
	go d.rpcS.ServeConn(d.server)

	b := make([]byte, BufferSize)
	for {
		n, err := d.origin.Read(b)
		if err != nil {
			return err
		}
		switch prefix, data := b[0], b[1:n]; prefix {
		case d.server.prefix:
			d.server.inCh <- data
		case d.client.prefix:
			d.client.inCh <- data
		default:
			return errors.New("invalid prefix")
		}
	}
}

func (d *RPCDuplex) Client() *rpc.Client {
	return d.rpcC
}
