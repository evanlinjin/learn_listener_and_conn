package main

import (
	"fmt"
	"log"
	"net"
)

// Writes a message 10 times to the echo server.
// Each message is prefixed with the message's index.
func main() {
	addr := "localhost:1234"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed to dial to address: %s", err)
	}
	log.Printf("dialed to address: %s\n", addr)

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("closed with error: %s", err)
			return
		}
		log.Println("closed successfully")
	}()

	msg := "This is the message to be sent!"

	for i := 0; i < 10; i++ {
		v := fmt.Sprintf("%d: %s", i+1, msg)

		// Write the message appended with the n-th number.
		n, err := conn.Write([]byte(v))
		if err != nil {
			log.Fatalf("failed to write to server: %s", err)
		}
		log.Printf("wrote: %s", v)

		// Read the response from the server.
		// We are expecting to read as much as we wrote as it is an echo server.
		resp := make([]byte, n)
		if _, err := conn.Read(resp); err != nil {
			log.Fatalf("failed to read from server: %s", err)
		}
		log.Printf("read: %s", string(resp))
	}
}

/* EXPECTED OUTPUT:

2019/04/27 16:41:52 dialed to address: localhost:1234
2019/04/27 16:41:52 wrote: 1: This is the message to be sent!
2019/04/27 16:41:52 read: 1: This is the message to be sent!
2019/04/27 16:41:52 wrote: 2: This is the message to be sent!
2019/04/27 16:41:52 read: 2: This is the message to be sent!
2019/04/27 16:41:52 wrote: 3: This is the message to be sent!
2019/04/27 16:41:52 read: 3: This is the message to be sent!
2019/04/27 16:41:52 wrote: 4: This is the message to be sent!
2019/04/27 16:41:52 read: 4: This is the message to be sent!
2019/04/27 16:41:52 wrote: 5: This is the message to be sent!
2019/04/27 16:41:52 read: 5: This is the message to be sent!
2019/04/27 16:41:52 wrote: 6: This is the message to be sent!
2019/04/27 16:41:52 read: 6: This is the message to be sent!
2019/04/27 16:41:52 wrote: 7: This is the message to be sent!
2019/04/27 16:41:52 read: 7: This is the message to be sent!
2019/04/27 16:41:52 wrote: 8: This is the message to be sent!
2019/04/27 16:41:52 read: 8: This is the message to be sent!
2019/04/27 16:41:52 wrote: 9: This is the message to be sent!
2019/04/27 16:41:52 read: 9: This is the message to be sent!
2019/04/27 16:41:52 wrote: 10: This is the message to be sent!
2019/04/27 16:41:52 read: 10: This is the message to be sent!
2019/04/27 16:41:52 closed successfully

*/