package configuration

import (
	"errors"
	"flag"
	"os"
)

type Configuration struct {
	ConfigurationFile string
	Chain             string
	Verbose           bool

	// the max port is 65535, a 16-bit number
	MetricsPort int
}

func GetFlagsConfiguration() Configuration {
	pConfigurationFile := flag.String("config", "./config.json", "Relative or absolute path to the JSON configuration file.")
	pChain := flag.String("chain", "INPUT", "Input chain name.")
	pMetricsPort := flag.Int("metrics-port", 9090, "Port used for the metrics server.")
	pVerbose := flag.Bool("v", false, "Enable DEBUG-level logging.")
	flag.Parse()

	return Configuration{
		ConfigurationFile: *pConfigurationFile,
		Chain:             *pChain,
		Verbose:           *pVerbose,
		MetricsPort:       *pMetricsPort,
	}
}

func CheckValues(c Configuration) error {
	if _, err := os.Open(c.ConfigurationFile); err != nil {
		return errors.New("Configuration file not found.")
	}

	if !(1 <= c.MetricsPort && c.MetricsPort <= 65535) {
		return errors.New("Invalid port for metrics, must be a value from 1 to 65535.")
	}

	// TODO: check if chain exists (?)

	return nil
}
