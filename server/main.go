package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type clientCh chan string

var (
	messages  = make(chan string)
	entringCh = make(chan clientCh)
	leavingCh = make(chan clientCh)
)

//var duration = time.Minute * 5

func broadcaster() {
	//var clients map[clientCh]bool
	clients := make(map[clientCh]string)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entringCh:
			clients[cli] = <-cli
			for clis := range clients {
				cli <- clients[clis]
			}
		case cli := <-leavingCh:
			delete(clients, cli)
			close(cli)
		}
	}

}
func handleConn(c net.Conn, name string) {
	duration := time.Minute * 5
	timer := time.NewTimer(duration)
	ch := make(chan string, 2)
	who := name //c.RemoteAddr().String()
	fmt.Println("WHO = ", who)
	entringCh <- ch
	ch <- who
	messages <- who + " connected"
	go connWriter(c, ch)

	ch <- "Вы :" + who
	input := bufio.NewScanner(c)
	go func(t **time.Timer) {
		for input.Scan() {
			messages <- who + ":" + input.Text()
			timer.Reset(duration)
		}
		<-(*timer).C
	}(&timer)
	<-timer.C
	fmt.Fprintln(c, "you disconnected")
	messages <- who + " disconnected"
	leavingCh <- ch
	c.Close()
}

//пишем в клиента все входящие сообщения
func connWriter(c net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(c, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err.Error())
	}
	go broadcaster()
	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		reader := bufio.NewReader(con)
		bytes := make([]byte, 20)
		//var bytes []byte
		n, err2 := reader.Read(bytes)
		if err2 != nil {
			fmt.Println(err.Error())
		} else {

			name := string(bytes[:n])
			go handleConn(con, name)
		}

	}
}
