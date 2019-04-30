package netutil

import (
	"github.com/stretchr/testify/require"
	"net"
	"net/rpc"
	"testing"
)

type RPC struct {}

func (RPC) Double(in *int, out *int) error {
	*out = *in * 2
	return nil
}

func TestRPCDuplex_Serve(t *testing.T) {
	cA, cB := net.Pipe()

	serverA := rpc.NewServer()
	require.NoError(t, serverA.RegisterName("RPC", new(RPC)))

	serverB := rpc.NewServer()
	require.NoError(t, serverB.RegisterName("RPC", new(RPC)))

	errChA := make(chan error)
	dA := NewRPCDuplex(cA, serverA, true)
	go func() {errChA <- dA.Serve()}()

	errChB := make(chan error)
	dB := NewRPCDuplex(cB, serverB, false)
	go func() {errChB <- dB.Serve()}()

	var r int
	for i := 0; i < 100; i++ {
		require.NoError(t, dA.Client().Call("RPC.Double", i, &r))
		require.Equal(t, i*2, r)

		require.NoError(t, dB.Client().Call("RPC.Double", i, &r))
		require.Equal(t, i*2, r)
	}
}
