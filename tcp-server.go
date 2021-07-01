package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	cmdAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6060")
	listener, err := net.ListenTCP("tcp", cmdAddr)

	var sm = semaphore.NewWeighted(int64(10))

	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	ctx := context.Background()
	fmt.Println("Server is listening...")

	wg := sync.WaitGroup{}
	for {
		select {
		case <-quitChan:
			listener.Close()
			wg.Wait()
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			continue
		}
		if err != nil {
			log.Println(err)
			continue
		}
		wg.Add(1)
		if err := sm.Acquire(ctx, 1); err != nil {
			listener.Close()
			wg.Wait()
			return
		}
		go func() {
			defer wg.Done()
			sm.Release(1)
			handleConnection(conn) //run the goroutine to process the request
		}()
		wg.Wait()
	}
}

// connection handling
func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		// read the data received in the request
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}
		source := string(input[0:n])
		target, err := strconv.Atoi(source)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		target *= target
		fmt.Println(source, "-", target)
		//send data to the client
		conn.Write([]byte(strconv.Itoa(target)))
	}
}
