package main

import (
	"fmt"
	"log"
	"net"
)

// =========================
// TASK
// =========================
// Write some example code with net.Listener
// listening on a given port (in a go routine).
// Read data from each net.Conn accepted and echo
// the same data back. Then, use net.Dial() to
// the listening port from the returned net.Conn.

func main() {

	servAddr := "localhost:1234"
	conn, err := net.Dial("tcp", servAddr)
	fmt.Printf("Client dialing to %s\n", servAddr)
	defer conn.Close()

	if err != nil {
		log.Fatal("Dial error", err)
	}

	// Writes to server
	fmt.Fprintln(conn, "Hello from client using fmt.Fprintln!")

	// This also writes to server but for some reason only executes after conn closes
	_, err = conn.Write([]byte("Hello using conn.Write([]byte)"))
	if err != nil {
		log.Fatal("Write to server failed", err)
	}

	// Read response from server
	// response, err := ioutil.ReadAll(conn)
	// if err != nil {
	// 	log.Fatal("Error reading response from server", err)
	// }
	// fmt.Println(string(response))

	// Read response from server
	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		log.Fatal("Error reading response from server", err)
	}

	fmt.Println("Server response", string(reply))

}
