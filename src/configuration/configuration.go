package configuration

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	ConfigurationFile string
	Chain             string
	CerberoSocket     string
	CerberoSocketIP   string
	CerberoSocketPort int
	Verbose           bool

	// the max port is 65535, a 16-bit number
	MetricsPort int
}

func GetFlagsConfiguration() Configuration {
	pConfigurationFile := flag.String("config-file", "", "Relative or absolute path to the JSON configuration file.")
	pChain := flag.String("chain", "INPUT", "Input chain name.")
	pCerberoSocket := flag.String("cerbero-socket", "", "The server to which Cerbero will connect to update the configuration file. The format must be <ip>:<port>.")
	pMetricsPort := flag.Int("metrics-port", 9090, "Port used for the metrics server.")
	pVerbose := flag.Bool("v", false, "Enable DEBUG-level logging.")
	flag.Parse()

	return Configuration{
		ConfigurationFile: *pConfigurationFile,
		Chain:             *pChain,
		Verbose:           *pVerbose,
		CerberoSocket:     *pCerberoSocket,
		MetricsPort:       *pMetricsPort,
	}
}

func CheckValues(c *Configuration) error {
	if c.ConfigurationFile == "" && c.CerberoSocket == "" {
		return errors.New("You must specify either --config-file or --cerbero-socket.")
	}

	if c.ConfigurationFile != "" {
		if _, err := os.Open(c.ConfigurationFile); err != nil {
			return errors.New("Configuration file not found.")
		}
	}

	if !(1 <= c.MetricsPort && c.MetricsPort <= 65535) {
		return errors.New("Invalid port for metrics, must be a value from 1 to 65535.")
	}

	// TODO: check if chain exists (?)

	if c.CerberoSocket != "" {
		cerberoSocketSplit := strings.Split(c.CerberoSocket, ":")
		if len(cerberoSocketSplit) != 2 {
			return errors.New("The Cerbero socket must be in the form of <ip>:<port>.")
		}

		var err error
		c.CerberoSocketIP = cerberoSocketSplit[0]
		c.CerberoSocketPort, err = strconv.Atoi(cerberoSocketSplit[1])
		if err != nil {
			return errors.New("The Cerbero socket port must be an integer.")
		}

		if !(1 <= c.CerberoSocketPort && c.CerberoSocketPort <= 65535) {
			return errors.New("Invalid port for Cerbero socket, must be a value from 1 to 65535.")
		}
	}

	return nil
}

func IsConfigFileSet(c Configuration) bool {
	return c.ConfigurationFile != ""
}

func IsCerberoSocketSet(c Configuration) bool {
	return c.CerberoSocket != "" &&
		c.CerberoSocketIP != "" &&
		(1 <= c.CerberoSocketPort && c.CerberoSocketPort <= 65535)
}