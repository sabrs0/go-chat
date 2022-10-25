package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func mustCopy(w io.Writer, r io.Reader, doneCh chan struct{}) {
	_, err := io.Copy(w, r)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	name := []byte(os.Args[1])
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		conn.Write(name)
		done := make(chan struct{})
		go func() {
			mustCopy(os.Stdout, conn, done)
			if err != nil {
				fmt.Println(err.Error())
			}

			done <- struct{}{}
		}()
		mustCopy(conn, os.Stdin, done)
		conn.Close()
		<-done
	}

}
