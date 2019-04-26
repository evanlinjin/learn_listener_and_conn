package main

import (
	"fmt"
	"io"
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

func handle(conn net.Conn) {

	// Reads data from net.Conn and echo all incoming data back
	// io.Copy(dst Writer, src Reader)
	// Basically reads from a reader and writes to a writer
	// defer conn.Close()
	io.Copy(conn, conn)

	// // Another echoer
	// var buf [512]byte
	// for {
	// 	// read upto 512 bytes
	// 	n, err := conn.Read(buf[0:])
	// 	if err != nil {
	// 		return
	// 	}

	// 	// write the n bytes read
	// 	_, err2 := conn.Write(buf[0:n])
	// 	if err2 != nil {
	// 		return
	// 	}
	// }

	// // Use bufio.NewScanner to read from connection
	// scanner := bufio.NewScanner(conn)
	// for scanner.Scan() {
	// 	ln := scanner.Text()
	// 	fmt.Println(ln)
	// 	fmt.Fprintf(conn, "Hello from server, please reply:")
	// }

}

func main() {

	// Listen on TCP port 1234
	port := ":1234"
	l, err := net.Listen("tcp", port)
	fmt.Printf("Server listen on port%s\n", port)

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		// Wait for a connection and accepts it
		// conn implements reader and writer
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Various methods to write from server to client using conn:
		io.WriteString(conn, "\nMethod1: io.WriteString(conn, string)\n")
		fmt.Fprintln(conn, "Method2: fmt.Fprintln(conn, string)")
		fmt.Fprintf(conn, "%v", "Method3: Fprintf(conn, string)")

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting so that
		// multiple connections may be served concurrently
		go handle(conn)

	}
}

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).
// func Copy(dst Writer, src Reader) (written int64, err error) {
// 	return copyBuffer(dst, src, nil)
// }

// func Copy(dst Writer, src Reader) (written int64, err error) {
// 	return copyBuffer(dst, src, nil)
// }

// https://godoc.org/io#Copy
// r := strings.NewReader("some io.Reader stream to be read\n")

// if _, err := io.Copy(os.Stdout, r); err != nil {
//     log.Fatal(err)
// }
// Output:

// some io.Reader stream to be read
