package configuration

import (
	"cerbero3/configuration/cli"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	ConfigurationFile string
	CerberoSocket     string
	CerberoSocketIP   string
	CerberoSocketPort int
	ColoredLogs       bool
	LogFile           string
	Verbose           bool

	// the max port is 65535, a 16-bit number
	MetricsPort int
}

func GetFlagsConfiguration() Configuration {
	structure := cli.Parse()
	return Configuration{
		ConfigurationFile: structure.File.FilePath,
		Verbose:           structure.Verbose,
		CerberoSocket:     structure.Socket.Address,
		ColoredLogs:       structure.ColoredLogs,
		LogFile:           structure.LogFile,
		MetricsPort:       structure.MetricsPort,
	}
}

func CheckValues(c *Configuration) error {
	if c.ConfigurationFile == "" && c.CerberoSocket == "" {
		return errors.New("You must specify either file or socket.")
	}

	if c.ConfigurationFile != "" {
		if _, err := os.Open(c.ConfigurationFile); err != nil {
			return errors.New("Configuration file not found.")
		}
	}

	if !(1 <= c.MetricsPort && c.MetricsPort <= 65535) {
		return errors.New("Invalid port for metrics, must be a value from 1 to 65535.")
	}

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
