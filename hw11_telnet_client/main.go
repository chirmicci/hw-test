package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connect to server timeout")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatalf("should set at 2 positional args")
	}

	client := NewTelnetClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer cancel()
		if err := client.Send(); err != nil {
			return
		}

		log.Print("...EOF")
	}()

	go func() {
		defer cancel()
		if err := client.Receive(); err != nil {
			return
		}
		log.Print("...Connection was closed by peer")
	}()

	<-ctx.Done()

	if err := client.Close(); err != nil {
		log.Printf("failed to close client: %v", err)
	}
}
