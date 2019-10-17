package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ut0mt8/goChecker/config"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	var wg sync.WaitGroup

	checks, err := config.GetConfig("checks.conf")
	if err != nil {
		log.Fatalf("configuration error : %v", err)
	}

	for _, c := range checks {
		wg.Add(1)
		go c.Start(done, &wg)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Info("beginning to serve on port 8080")
	go http.ListenAndServe(":8080", nil)

	<-quit
	close(done)
	log.Info("stoping. waiting for probes to finish...")
	wg.Wait()
	log.Info("stop.")

}
