package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"svcproxy/config"
	"svcproxy/service"
)

// Version to be filled by ldflags
var Version = "dev"

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/svcproxy/services.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Error parsing configuration: %s", err)
	}

	svc, err := service.NewService()
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	for _, sd := range cfg.Services {
		backend, err := url.Parse(sd.Backend.URL)
		if err != nil {
			log.Fatalf("Error parsing url: %s", err)
		}

		svc.AddProxy(&service.Proxy{
			Frontend: &service.Frontend{
				FQDN: sd.Frontend.FQDN,
			},
			Backend: &service.Backend{
				URL: backend,
			},
		})
	}

	httpSvc := &http.Server{
		Addr:    cfg.Listener.HTTPAddr,
		Handler: svc,
	}
	httpsSvc := &http.Server{
		Addr:    cfg.Listener.HTTPSAddr,
		Handler: svc,
	}

	go func() { log.Fatal(httpSvc.ListenAndServe()) }()
	log.Fatal(httpsSvc.ListenAndServe())
}
