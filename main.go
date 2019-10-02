package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ut0mt8/goChecker/config"
)

func main() {
	checks, err := config.GetConfig("checks.conf")
	if err != nil {
		log.Fatalf("configuration error : %v", err)
	}

	for _, c := range checks {
		go c.Start()
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Info("beginning to serve on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
