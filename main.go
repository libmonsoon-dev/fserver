package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port      int
	directory string
)

func main() {
	flag.IntVar(&port, "port", 0, "listen port")
	flag.StringVar(&directory, "dir", ".", "directory")
	flag.Parse()

	ctx, stopNotify := signal.NotifyContext(context.Background(), signals...)
	defer stopNotify()

	localAddr := &net.TCPAddr{Port: port}
	listener := must(net.ListenTCP("tcp4", localAddr))
	defer listener.Close()

	printLocalAddress(listener.Addr().(*net.TCPAddr).Port)

	server := http.Server{
		BaseContext: func(ln net.Listener) context.Context {
			return ctx
		},
		Handler: http.FileServer(http.Dir(directory)),
	}
	defer server.Close()

	go server.Serve(listener)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func printLocalAddress(port int) {
	log.Printf("Listening on:")
	for _, iface := range must(net.Interfaces()) {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		for _, addr := range must(iface.Addrs()) {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}

			log.Printf("\thttp://%s/", &net.TCPAddr{IP: ip, Port: port})
		}
	}
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

var signals = []os.Signal{
	syscall.SIGTERM,
	syscall.SIGINT,
	syscall.SIGHUP,
}
