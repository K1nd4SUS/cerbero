package metrics

import (
	"cerbero/configuration"
	"cerbero/logs"
	"cerbero/metrics/metricsdb"
	"cerbero/metrics/metricsregex"
	"cerbero/services"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func IncrementService(serviceIndex int, dropped bool) {
	serviceCounters := metricsdb.GetServiceCounters(services.Services[serviceIndex])
	serviceCounters.TotalPackets.Inc()
	if dropped {
		serviceCounters.DroppedPackets.Inc()
	}

	logs.PrintDebug(fmt.Sprintf("Incremented prometheus counter for service %v (total%v).", services.Services[serviceIndex].Name, func() string {
		if dropped {
			return ", dropped"
		} else {
			return ""
		}
	}()))
}

func IncrementRegex(serviceIndex int, regex string) {
	regexCounter := metricsdb.GetRegexCounter(services.Services[serviceIndex], regex)
	regexCounter.Inc()

	logs.PrintDebug(fmt.Sprintf("Incremented prometheus counter for regex %v (dropped) in service %v.", metricsregex.ToHex(regex), services.Services[serviceIndex].Name))
}

func StartServer(config configuration.Configuration) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.MetricsPort), nil)
	if err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
}
