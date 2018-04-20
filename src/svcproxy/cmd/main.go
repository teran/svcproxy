package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"runtime"

	"golang.org/x/crypto/acme/autocert"

	"svcproxy/authentication/factory"
	"svcproxy/autocert/cache"
	"svcproxy/config"
	"svcproxy/middleware"
	"svcproxy/service"
)

// Version to be filled by ldflags
var Version = "dev"

func main() {
	log.Printf("Launching svcproxy=%s ...", Version)
	log.Printf("Built with %s", runtime.Version())

	// Grab path to configuration file and load it
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/svcproxy/svcproxy.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Error parsing configuration: %s", err)
	}

	// Create service instance
	svc, err := service.NewService()
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	var hostsList []string

	// Fill service instance with proxies
	for _, sd := range cfg.Services {
		for _, fqdn := range sd.Frontend.FQDN {
			f, err := service.NewFrontend(fqdn, sd.Frontend.HTTPHandler, sd.Frontend.ResponseHTTPHeaders)
			if err != nil {
				log.Printf("Error: unable to initialize frontend %s: %s", fqdn, err)
				continue
			}

			a, err := factory.NewAuthenticator(sd.Authentication.Method, sd.Authentication.Options)
			if err != nil {
				log.Printf("Error: unable to initialize authenticator %s: %s", sd.Authentication.Method, err)
				log.Printf("Skipping service %s from initialization", fqdn)
				continue
			}

			b, err := service.NewBackend(sd.Backend.URL)
			if err != nil {
				log.Printf("Error: unable to initialize backend: %s: %s", fqdn, err)
				log.Printf("Skipping service %s from initialization", fqdn)
				continue
			}

			p, err := service.NewProxy(f, b, a)
			if err != nil {
				log.Printf("Error: unable to register proxy %s: %s", fqdn, err)
				log.Printf("Skipping service %s from initialization", fqdn)
				continue
			}
			svc.AddProxy(p)

			hostsList = append(hostsList, fqdn)
		}
	}

	// Initialize caching subsystem
	cache, err := cache.NewCacheFactory(cfg.Autocert.Cache.Backend, cfg.Autocert.Cache.BackendOptions)
	if err != nil {
		log.Fatalf("Error initializing autocert cache: %s", err)
	}

	log.Print("Loaded proxies for hosts:")
	for _, host := range hostsList {
		log.Printf(" - %s", host)
	}

	// Initialize autocert
	acm := &autocert.Manager{
		Cache:      cache,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostsList...),
	}

	debugSvc := &http.Server{
		Addr:    cfg.Listener.DebugAddr,
		Handler: http.HandlerFunc(svc.DebugHandlerFunc),
	}
	go func() {
		log.Printf("Listening to Debug HTTP socket: %s", cfg.Listener.DebugAddr)
		log.Fatalf("Error listening Debug HTTP socket: %s", debugSvc.ListenAndServe())
	}()

	// Run http listeners
	httpSvc := &http.Server{
		Addr:    cfg.Listener.HTTPAddr,
		Handler: acm.HTTPHandler(svc),
	}
	go func() {
		log.Printf("Listening to Service HTTP socket: %s", cfg.Listener.HTTPAddr)
		log.Fatalf("Error listening Service HTTP socket: %s", httpSvc.ListenAndServe())
	}()

	// Configure TLS
	tlsconf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		},
		GetCertificate:           acm.GetCertificate,
		PreferServerCipherSuites: true,
	}

	// Run HTTPS listener
	httpsSvc := &http.Server{
		Addr:      cfg.Listener.HTTPSAddr,
		TLSConfig: tlsconf,
		Handler:   middleware.Chain(svc, cfg.Listener.Middlewares...),
	}
	log.Printf("Listening to Service HTTPS socket: %s", cfg.Listener.HTTPSAddr)
	log.Fatalf("Error listening HTTPS socket: %s", httpsSvc.ListenAndServeTLS("", ""))
}
