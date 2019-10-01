package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ut0mt8/goChecker/config"
	"net/http"
)

func main() {
	checks, err := config.GetConfig("checks.conf")
	if err != nil {
		log.Fatalf("configuration error : %v", err)
	}

	for _, c := range checks.Check {
		go c.Start()
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Info("beginning to serve on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
