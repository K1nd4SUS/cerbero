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
}

func GetFlagsConfiguration() Configuration {
	pConfigurationFile := flag.String("config", "./config.json", "Relative or absolute path to the JSON configuration file.")
	pChain := flag.String("chain", "INPUT", "Input chain name.")
	pVerbose := flag.Bool("v", false, "Enable DEBUG-level logging.")
	flag.Parse()

	return Configuration{
		ConfigurationFile: *pConfigurationFile,
		Chain:             *pChain,
		Verbose:           *pVerbose,
	}
}

func CheckValues(c Configuration) error {
	if _, err := os.Open(c.ConfigurationFile); err != nil {
		return errors.New("Configuration file not found.")
	}

	// TODO: check if chain exists (?)

	return nil
}
