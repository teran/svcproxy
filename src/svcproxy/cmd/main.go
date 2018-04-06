package main

import (
	"fmt"
	// svc "svcproxy/service"
)

type config struct {
	HTTPAddr  string `default:":80"`
	HTTPSAddr string `default:":443"`
}

// Version to be filled by ldflags
var Version = "dev"

func main() {
	fmt.Printf("Running svcproxy=%s\n", Version)
}
