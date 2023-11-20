package services

import (
	"bytes"
	"cerbero3/logs"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var Services []Service

type Service struct {
	Nfq       uint16   `json:"nfq"`
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Port      int      `json:"port"`
	RegexList []string `json:"regex_list"`
}

func Load(configurationFile string) error {
	logs.PrintDebug(fmt.Sprintf(`Reading configuration file from "%v"...`, configurationFile))
	b, err := os.ReadFile(configurationFile)
	if err != nil {
		return err
	}
	logs.PrintDebug(fmt.Sprintf(`Read configuration file from "%v".`, configurationFile))

	logs.PrintDebug("Parsing JSON file...")
	err = json.Unmarshal(b, &Services)
	if err != nil {
		return err
	}
	logs.PrintDebug(func() string {
		buffer := &bytes.Buffer{}
		// TODO: this should have been already parsed by the Unmarshal
		// function, if it outputs an error there's clearly something
		// wrong with the compiler or something else
		json.Compact(buffer, b)

		return fmt.Sprintf("Parsed JSON file: %v", buffer.String())
	}())

	return nil
}

func CheckServicesValues() error {
	for _, service := range Services {
		if !(1 <= service.Nfq && service.Nfq <= 65535) {
			return errors.New(fmt.Sprintf(`Invalid nfq for service %v, must be a value from 1 to 65535.`, service.Name))
		}

		if !(service.Protocol == "tcp" || service.Protocol == "udp") {
			return errors.New(fmt.Sprintf(`Invalid protocol for service %v, must be set to either "tcp" or "udp".`, service.Name))
		}

		if !(1 <= service.Port && service.Port <= 65535) {
			return errors.New(fmt.Sprintf(`Invalid port for service %v, must be a value from 1 to 65535.`, service.Name))
		}
	}

	return nil
}
