package main

import (
	"log"
	"net/http"
	"net/url"

	"svcproxy/service"
)

type config struct {
	HTTPAddr  string `default:":80"`
	HTTPSAddr string `default:":443"`
}

// Version to be filled by ldflags
var Version = "dev"

func main() {
	frontend := "localhost"
	backend, err := url.Parse("http://ya.ru")
	if err != nil {
		log.Fatalf("Error parsing url: %s", err)
	}

	svc, err := service.NewService()
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	svc.AddProxy(&service.Proxy{
		Frontend: &service.Frontend{
			FQDN: frontend,
		},
		Backend: &service.Backend{
			URL: backend,
		},
	})

	http.ListenAndServe(":8080", svc)
}
