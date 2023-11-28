package metrics

import (
	"cerbero3/logs"
	"cerbero3/metrics/metricsdb"
	"cerbero3/metrics/metricsregex"
	"cerbero3/services"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func IncrementService(service services.Service, dropped bool) {
	serviceCounters := metricsdb.GetServiceCounters(service)
	serviceCounters.TotalPackets.Inc()
	if dropped {
		serviceCounters.DroppedPackets.Inc()
	}

	logs.PrintDebug(fmt.Sprintf("Incremented prometheus counter for service %v (total%v).", service.Name, func() string {
		if dropped {
			return ", dropped"
		} else {
			return ""
		}
	}()))
}

func IncrementRegex(regex string) {
	regexCounter := metricsdb.GetRegexCounter(regex)
	regexCounter.Inc()

	logs.PrintDebug(fmt.Sprintf("Incremented prometheus counter for regex %v (dropped).", metricsregex.ToHex(regex)))
}

func StartServer() {
	// TODO: make the port configurable
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
}
