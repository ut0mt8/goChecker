package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ut0mt8/goChecker/checker"
	"github.com/ut0mt8/goChecker/config"
	"net/http"
)

func main() {
	checks, err := config.GetConfig("checks.conf")
	if err != nil {
		log.Fatalf("error reading configuration : %v", err)
	}

	for _, c := range checks.Check {
		go checker.StartCheck(c)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}