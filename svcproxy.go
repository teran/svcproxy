package main

import (
	"crypto/tls"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"runtime"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	log "github.com/sirupsen/logrus"
	"github.com/teran/svcproxy/authentication/factory"
	"github.com/teran/svcproxy/autocert/cache"
	"github.com/teran/svcproxy/config"
	"github.com/teran/svcproxy/middleware"
	"github.com/teran/svcproxy/service"
)

// Version to be filled by ldflags
var Version = "dev"

func main() {
	// Grab path to configuration file and load it
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/svcproxy/svcproxy.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.WithFields(log.Fields{
			"reason": err,
		}).Fatal("Error parsing configuration")
	}

	setLogFormatter(cfg.Logger.Formatter)
	setLogLevel(cfg.Logger.Level)
	logStartupHeader()

	// Create service instance
	svc, err := service.NewService()
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	var logger = log.New()
	w := logger.Writer()
	defer w.Close()

	var hostsList []string

	// Fill service instance with proxies
	for _, sd := range cfg.Services {
		for _, fqdn := range sd.Frontend.FQDN {
			f, err := service.NewFrontend(fqdn, sd.Frontend.HTTPHandler, sd.Frontend.ResponseHTTPHeaders)
			if err != nil {
				log.WithFields(log.Fields{
					"reason": err,
					"object": fqdn,
				}).Warn("Error: unable to initialize frontend. Skipping.")
				continue
			}

			a, err := factory.NewAuthenticator(sd.Authentication.Method, sd.Authentication.Options)
			if err != nil {
				log.WithFields(log.Fields{
					"reason": err,
					"object": sd.Authentication.Method,
					"parent": fqdn,
				}).Warn("Error: unable to initialize auhenticator. Skipping.")
				continue
			}

			b, err := service.NewBackend(sd.Backend.URL, sd.Backend.RequestHTTPHeaders)
			if err != nil {
				log.WithFields(log.Fields{
					"reason": err,
					"object": sd.Backend.URL,
					"parent": fqdn,
				}).Warn("Error: unable to initialize backend. Skipping.")
				continue
			}

			transport := &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   cfg.Listener.Backend.Timeout,
					KeepAlive: cfg.Listener.Backend.KeepAlive,
					DualStack: cfg.Listener.Backend.DualStack,
				}).DialContext,
				ExpectContinueTimeout: cfg.Listener.Backend.ExpectContinueTimeout,
				IdleConnTimeout:       cfg.Listener.Backend.IdleConnTimeout,
				MaxIdleConns:          cfg.Listener.Backend.MaxIdleConns,
				Proxy:                 http.ProxyFromEnvironment,
				ResponseHeaderTimeout: cfg.Listener.Backend.ResponseHeaderTimeout,
				TLSHandshakeTimeout:   cfg.Listener.Backend.TLSHandshakeTimeout,
			}

			p, err := service.NewProxy(f, b, a, transport, stdlog.New(w, "", 0))
			if err != nil {
				log.WithFields(log.Fields{
					"reason": err,
					"object": fqdn,
				}).Warn("Error: unable to register proxy. Skipping.")
				continue
			}
			svc.AddProxy(p)

			hostsList = append(hostsList, fqdn)
		}
	}

	cache := initializeCache(cfg.Autocert.Cache.Backend, cfg.Autocert.Cache.BackendOptions)

	log.Debug("Loaded proxies for hosts:")
	for _, host := range hostsList {
		log.Debugf(" - %s", host)
	}

	// Initialize autocert
	acm := &autocert.Manager{
		Email:      cfg.Autocert.Email,
		Cache:      cache,
		Client:     &acme.Client{DirectoryURL: cfg.Autocert.DirectoryURL},
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostsList...),
	}

	debugSvc := &http.Server{
		Addr:    cfg.Listener.DebugAddr,
		Handler: http.HandlerFunc(svc.DebugHandlerFunc),
	}
	go func() {
		log.WithFields(log.Fields{
			"socket": cfg.Listener.DebugAddr,
		}).Info("Listening to Debug HTTP socket")

		err = debugSvc.ListenAndServe()
		log.WithFields(log.Fields{
			"reason": err,
		}).Fatal("Error listening Debug HTTP socket")
	}()

	// Run http listeners
	httpSvc := &http.Server{
		Addr:              cfg.Listener.HTTPAddr,
		Handler:           middleware.Chain(acm.HTTPHandler(svc), cfg.Listener.Middlewares...),
		IdleTimeout:       cfg.Listener.Frontend.IdleTimeout,
		ReadHeaderTimeout: cfg.Listener.Frontend.ReadHeaderTimeout,
		ReadTimeout:       cfg.Listener.Frontend.ReadTimeout,
		WriteTimeout:      cfg.Listener.Frontend.WriteTimeout,
	}
	go func() {
		log.WithFields(log.Fields{
			"socket": cfg.Listener.HTTPAddr,
		}).Info("Listening to Service HTTP socket")

		err = httpSvc.ListenAndServe()
		log.WithFields(log.Fields{
			"reason": err,
		}).Fatal("Error listening Service HTTP socket")
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
		Addr:              cfg.Listener.HTTPSAddr,
		Handler:           middleware.Chain(svc, cfg.Listener.Middlewares...),
		TLSConfig:         tlsconf,
		IdleTimeout:       cfg.Listener.Frontend.IdleTimeout,
		ReadHeaderTimeout: cfg.Listener.Frontend.ReadHeaderTimeout,
		ReadTimeout:       cfg.Listener.Frontend.ReadTimeout,
		WriteTimeout:      cfg.Listener.Frontend.WriteTimeout,
	}
	log.WithFields(log.Fields{
		"socket": cfg.Listener.HTTPSAddr,
	}).Info("Listening to Service HTTPS socket")

	err = httpsSvc.ListenAndServeTLS("", "")
	log.WithFields(log.Fields{
		"reason": err,
	}).Fatal("Error listening HTTPS socket")
}

func setLogFormatter(formatter string) {
	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.WithFields(log.Fields{
			"reason": fmt.Sprintf("unknown formatter '%s'", formatter),
		}).Fatalf("Error configuring logger")
	}
}

func initializeCache(backend string, options map[string]string) autocert.Cache {
	// Initialize caching subsystem
	cache, err := cache.NewCacheFactory(backend, options)
	if err != nil {
		log.WithFields(log.Fields{
			"reason": err,
		}).Fatal("Error: unable to initialize autocert cache")
	}

	return cache
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.WithFields(log.Fields{
			"reason": fmt.Sprintf("unknown log level '%s'", level),
		}).Fatalf("Error configuring logging level")
	}
}

func logStartupHeader() {
	log.WithFields(log.Fields{
		"version": Version,
	}).Infof("Launching svcproxy...")
	log.Debugf("Built with %s", runtime.Version())
}
