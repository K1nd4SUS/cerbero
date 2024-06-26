package metricsdb

import (
	"cerbero/metrics/metricsjobs"
	"cerbero/metrics/metricsregex"
	"cerbero/services"
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
// second key is the regex
// value is the counter(s)
var regexesDatabase = make(map[string]map[string]prometheus.Counter)

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
		if serviceCountersJob.IsActive() {
			serviceCountersJob.Wait()
			return GetServiceCounters(service)
		}

		return createServiceCounters(service)
	}
	return serviceCounters
}

func createRegexCounter(service services.Service, regex string) prometheus.Counter {
	regexCountersJob.Add(1)

	newRegexCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("cerbero_regex_%v_%v_packets_dropped", service.Name, metricsregex.ToHex(regex)),
		Help: fmt.Sprintf(`The number of packets that were blocked from regex "%v" (hex).`, metricsregex.ToHex(regex)),
	})

	// we made a map of maps, but of course
	// if we need to assign a value to a second-level
	// map, we need to make sure that the first-level
	// map exists first
	_, ok := regexesDatabase[service.Name]
	if !ok {
		regexesDatabase[service.Name] = make(map[string]prometheus.Counter)
	}
	regexesDatabase[service.Name][regex] = newRegexCounter

	regexCountersJob.Done()
	return newRegexCounter
}

func GetRegexCounter(service services.Service, regex string) prometheus.Counter {
	regexCounter, ok := regexesDatabase[service.Name][regex]

	// if regex already exists as a key in the servicesDatabase,
	// then "ok" is going to be true and it's going to return it;
	// else, it's going to create a new entry
	if !ok {
		if regexCountersJob.IsActive() {
			regexCountersJob.Wait()
			return GetRegexCounter(service, regex)
		}

		return createRegexCounter(service, regex)
	}
	return regexCounter
}
