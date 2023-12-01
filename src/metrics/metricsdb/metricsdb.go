package metricsdb

import (
	"cerbero3/metrics/metricsjobs"
	"cerbero3/metrics/metricsregex"
	"cerbero3/services"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type serviceCounters struct {
	TotalPackets   prometheus.Counter
	DroppedPackets prometheus.Counter
}

// key is the service name
// value is the counter(s)
var servicesDatabase = make(map[string]serviceCounters)

// key is the service name
// value is the counter(s)
var regexesDatabase = make(map[string]prometheus.Counter)

var serviceCountersJob = metricsjobs.CreatingCountersJob{}
var regexCountersJob = metricsjobs.CreatingCountersJob{}

func createServiceCounters(service services.Service) serviceCounters {
	serviceCountersJob.Add(1)

	newServiceCounters := serviceCounters{
		TotalPackets: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("cerbero_service_%v_packets_total", service.Name),
			Help: fmt.Sprintf("The total number of packets that passed through port %v.", service.Port),
		}),
		DroppedPackets: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("cerbero_service_%v_packets_dropped", service.Name),
			Help: fmt.Sprintf("The number of packets that were blocked from port %v.", service.Port),
		}),
	}
	servicesDatabase[service.Name] = newServiceCounters

	serviceCountersJob.Done()
	return newServiceCounters
}

func GetServiceCounters(service services.Service) serviceCounters {
	serviceCounters, ok := servicesDatabase[service.Name]

	// if service.Name already exists as a key in the servicesDatabase,
	// then "ok" is going to be true and it's going to return it;
	// else, it's going to create a new entry
	if !ok {
		// if serviceCountersJob.IsActive() {
		// 	serviceCountersJob.Wait()
		// 	return GetServiceCounters(service)
		// }

		return createServiceCounters(service)
	}
	return serviceCounters
}

func createRegexCounter(regex string) prometheus.Counter {
	regexCountersJob.Add(1)

	newRegexCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("cerbero_regex_%v_packets_dropped", metricsregex.ToHex(regex)),
		Help: fmt.Sprintf(`The number of packets that were blocked from regex "%v" (hex).`, metricsregex.ToHex(regex)),
	})
	regexesDatabase[regex] = newRegexCounter

	regexCountersJob.Done()
	return newRegexCounter
}

func GetRegexCounter(regex string) prometheus.Counter {
	regexCounter, ok := regexesDatabase[regex]

	// if regex already exists as a key in the servicesDatabase,
	// then "ok" is going to be true and it's going to return it;
	// else, it's going to create a new entry
	if !ok {
		if regexCountersJob.IsActive() {
			regexCountersJob.Wait()
			return GetRegexCounter(regex)
		}

		return createRegexCounter(regex)
	}
	return regexCounter
}
