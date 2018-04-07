package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"

	sqlcache "svcproxy/autocert/cache/sql"
	"svcproxy/config"
	"svcproxy/service"

	"golang.org/x/crypto/acme/autocert"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Version to be filled by ldflags
var Version = "dev"

func main() {
	// Grab path to configuration file and load it
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/svcproxy/services.yaml"
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

		hostsList = append(hostsList, sd.Frontend.FQDN)
	}

	log.Printf("Loaded hosts: %+s", hostsList)

	// Initialize database for caching TLS certificates
	db, err := sql.Open("mysql", cfg.Autocert.Cache.BackendOptions["dsn"])
	if err != nil {
		log.Fatalf("Error establising database connection: %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while sending ping to database: %s", err)
	}

	cache, err := sqlcache.NewCache(db, []byte(cfg.Autocert.Cache.BackendOptions["encryptionKey"]))
	if err != nil {
		log.Fatalf("Error initializing autocert cache: %s", err)
	}

	acm := &autocert.Manager{
		Cache:      cache,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostsList...),
	}

	httpSvc := &http.Server{
		Addr:    cfg.Listener.HTTPAddr,
		Handler: acm.HTTPHandler(svc),
	}
	go func() {
		log.Fatalf("Error listening HTTP socket: %s", httpSvc.ListenAndServe())
	}()

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

	httpsSvc := &http.Server{
		Addr:      cfg.Listener.HTTPSAddr,
		TLSConfig: tlsconf,
		Handler:   svc,
	}
	log.Fatalf("Error listening HTTPS socket: %s", httpsSvc.ListenAndServeTLS("", ""))
}
