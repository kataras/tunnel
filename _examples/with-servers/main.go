package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/tunnel"
)

func main() {
	srv1 := &http.Server{Addr: ":8080", Handler: handler("server 1")}
	go srv1.ListenAndServe()

	srv2 := &http.Server{Addr: ":9090", Handler: handler("server 2")}
	go srv2.ListenAndServe()

	publicAddr := tunnel.MustStart(tunnel.WithServers(srv1, srv2))
	/*
		A shortcut of:
		tunnel.MustStart(tunnel.Configuration{
			Tunnels: []tunnel.Tunnel{
				{Addr: ":8080"},
				{Addr: ":9090"},
			},
		})
		And RegisterOnShutdown a tunnel.StopTunnel for each http server.
	*/
	fmt.Printf("• Public Addresses: %s\n", publicAddr)

	// Wait for kill -SIGINT XXXX or Ctrl/CMD+C
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-ch
}

func startServer(addr string) error {
	router := http.NewServeMux()
	router.HandleFunc("/", handler(addr))
	srv := &http.Server{Addr: addr, Handler: router}
	return srv.ListenAndServe()
}

func handler(srvName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from server: %s\n", srvName)
	}
}
