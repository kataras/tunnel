package main

import (
	"fmt"
	"net/http"

	"github.com/kataras/tunnel"
)

const addr = ":8080"

func main() {
	http.HandleFunc("/", handler)

	go func() {
		// Start tunnel.
		config := tunnel.Configuration{
			Tunnels: []tunnel.Tunnel{
				{Addr: addr},
			},
		}
		publicAddrs := tunnel.MustStart(config)
		// That's all.
		fmt.Printf("• Public Address: %s\n", publicAddrs[0])
	}()

	http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from server")
}
