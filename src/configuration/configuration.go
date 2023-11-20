package configuration

import (
	"errors"
	"flag"
	"os"
)

type Configuration struct {
	// TODO: check if this flag needs to be used or not
	// Nfq               int
	ConfigurationFile string
	Chain             string
	Verbose           bool
}

func GetFlagsConfiguration() Configuration {
	// pNfq := flag.Int("nfq", 100, "Queue number from 1 through 65535 (default: 100).")
	pConfigurationFile := flag.String("config", "./config.json", "Relative or absolute path to the JSON configuration file.")
	pChain := flag.String("chain", "INPUT", "Input chain name.")
	pVerbose := flag.Bool("v", false, "Enable DEBUG-level logging.")
	flag.Parse()

	return Configuration{
		// Nfq:               *pNfq,
		ConfigurationFile: *pConfigurationFile,
		Chain:             *pChain,
		Verbose:           *pVerbose,
	}
}

func CheckValues(c Configuration) error {
	// if !(1 <= c.Nfq && c.Nfq <= 65535) {
	// 	return errors.New("Invalid argument for flag -nfq, the value needs to sit between 1 and 65535.")
	// }

	if _, err := os.Open(c.ConfigurationFile); err != nil {
		return errors.New("Configuration file not found.")
	}

	// TODO: check if chain exists (?)

	return nil
}
