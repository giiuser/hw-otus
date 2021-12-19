package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Invalid arguments - host and port not define")
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		defer cancel()
		if err := client.Send(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer cancel()
		if err := client.Receive(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-ctx.Done()
}
