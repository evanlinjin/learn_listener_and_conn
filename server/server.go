package main

import (
	"io"
	"log"
	"net"
)

// Listens for connections and acts as an echo server for each connection.
func main() {
	addr := ":1234"

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	log.Printf("listening on address: %s", addr)

	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("closed with error: %s", err)
			return
		}
		log.Println("closed successfully")
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	// The following replicates io.Copy()
	for {
		// assumes we won't read data longer than 1024 bytes.
		b := make([]byte, 1024)

		// Read from conn.
		n, err := conn.Read(b)
		if err != nil {
			switch err {
			case io.EOF:
				log.Printf("connection %s has ended!", conn.RemoteAddr())
			default:
				log.Printf("failed to read from %s: %s", conn.RemoteAddr(), err)
			}
			return
		}
		log.Printf("read %d bytes from %s!", n, conn.RemoteAddr())

		// Echo the read bytes back to conn.
		if _, err := conn.Write(b[:n]); err != nil {
			log.Printf("failed to write to %s: %s", conn.RemoteAddr(), err)
			return
		}
		log.Printf("wrote %d bytes to %s!", n, conn.RemoteAddr())
	}
}

/* EXPECTED OUTPUT:

2019/04/27 16:41:48 listening on address: :1234
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 34 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 34 bytes to [::1]:50807!
2019/04/27 16:41:52 read 35 bytes from [::1]:50807!
2019/04/27 16:41:52 wrote 35 bytes to [::1]:50807!
2019/04/27 16:41:52 connection [::1]:50807 has ended!

*/