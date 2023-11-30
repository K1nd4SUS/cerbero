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
	MetricsPort int16
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
		MetricsPort:       int16(*pMetricsPort),
	}
}

func CheckValues(c Configuration) error {
	if _, err := os.Open(c.ConfigurationFile); err != nil {
		return errors.New("Configuration file not found.")
	}

	// TODO: check if chain exists (?)

	return nil
}
