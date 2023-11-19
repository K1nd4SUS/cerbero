package main

import (
        "net/http"
        "time"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
        go func() {
                for {
                        service.Inc()
                        time.Sleep(1 * time.Second)
                }
        }()
}

var (
        service = promauto.NewCounter(prometheus.CounterOpts{
                Name: "cerbero_8080_packets_total",
                Help: "The total number of packets that arriving on port 8080",
        })
)

func main() {
        recordMetrics()

        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}