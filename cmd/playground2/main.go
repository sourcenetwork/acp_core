// Command playground2 serves the Playground2 API over Connect unary JSON.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sourcenetwork/acp_core/pkg/services"
)

func main() {
	addr := flag.String("addr", "localhost:8090", "listen address")
	alloy := flag.String("alloy", "", "path to the alloy executable (default: look up in PATH)")
	flag.Parse()

	h := NewHandler(services.NewPlayground2Service(*alloy))
	log.Printf("playground2 listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, h))
}
