package netutil

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/rpc"
)

const CacheCount = 100

type branchConn struct {
	net.Conn             // original conn.
	inCh     chan []byte // original conn pushes here.
	prefix   byte        // Prefix to append before sending to original conn.
}

func newBranchConn(origin net.Conn, prefix byte) *branchConn {
	return &branchConn{Conn: origin, inCh: make(chan []byte, CacheCount), prefix: prefix}
}

func (s *branchConn) Read(p []byte) (n int, err error) {
	data, ok := <-s.inCh
	if !ok {
		return 0, io.ErrClosedPipe
	}
	if len(data) > len(p) {
		return 0, io.ErrShortBuffer
	}
	copy(p, data)
	return len(data), nil
}

func (s *branchConn) Write(p []byte) (n int, err error) {
	h := make([]byte, 3)
	h[0] = s.prefix
	binary.BigEndian.PutUint16(h[1:], uint16(len(p)))

	n, err = s.Conn.Write(append(h, p...))
	return n - 3, err
}

type RPCDuplex struct {
	connO net.Conn    // original conn
	connS *branchConn // server conn
	connC *branchConn // client conn
	rpcS  *rpc.Server // rpc server
	rpcC  *rpc.Client // rpc client
}

// NewRPCDuplex creates a new RPCDuplex with a given connection and rpc server.
// 'init' specifies whether this instance is the initiator.
func NewRPCDuplex(conn net.Conn, srv *rpc.Server, init bool) *RPCDuplex {
	var serverPrefix, clientPrefix byte
	if init {
		serverPrefix, clientPrefix = 0, 1
	} else {
		serverPrefix, clientPrefix = 1, 0
	}
	serverConn := newBranchConn(conn, serverPrefix)
	clientConn := newBranchConn(conn, clientPrefix)
	return &RPCDuplex{connO: conn, connS: serverConn, connC: clientConn, rpcS: srv, rpcC: rpc.NewClient(clientConn)}
}

// Serve serves the RPC server and runs the event loop that forwards data to branchConns.
func (d *RPCDuplex) Serve() error {
	go d.rpcS.ServeConn(d.connS)

	for {
		h := make([]byte, 3)
		if _, err := io.ReadFull(d.connO, h); err != nil {
			return err
		}
		prefix := h[0]
		size := binary.BigEndian.Uint16(h[1:])

		p := make([]byte, size)
		if _, err := io.ReadFull(d.connO, p); err != nil {
			return err
		}
		switch prefix {
		case d.connS.prefix:
			d.connS.inCh <- p
		case d.connC.prefix:
			d.connC.inCh <- p
		default:
			return errors.New("invalid prefix")
		}
	}
}

// Client returns the internal RPC Client.
func (d *RPCDuplex) Client() *rpc.Client { return d.rpcC }
