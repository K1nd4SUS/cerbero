package rules

import (
	"cerbero/configuration"
	"cerbero/logs"
	"cerbero/services"
	"fmt"
	"os"
	"os/exec"
)

type RemoveRules chan bool

func getIPTablesCommand(config configuration.Configuration, service services.Service, doAdd bool) *exec.Cmd {
	return exec.Command(
		"iptables",
		// if doAdd is true, it is "-I", else it is "-D"
		fmt.Sprintf("-%v", func() string {
			if doAdd {
				return "I"
			} else {
				return "D"
			}
		}()), service.Chain,
		"-p", service.Protocol,
		// taken from:
		// https://docs.docker.com/network/packet-filtering-firewalls/#match-the-original-ip-and-ports-for-requests
		"-m", "conntrack",
		"--ctorigdstport", fmt.Sprintf("%v", service.Port),
		"-j", "NFQUEUE",
		"--queue-num", fmt.Sprintf("%v", service.Nfq),
	)
}

func SetRules(config configuration.Configuration) error {
	for _, service := range services.Services {
		cmd := getIPTablesCommand(config, service, true)

		logs.PrintDebug(fmt.Sprintf(`Setting rules for service "%v" with command "%v"...`, service.Name, cmd.String()))
		err := cmd.Run()
		if err != nil {
			return err
		}
		logs.PrintInfo(fmt.Sprintf(`Set rules for service "%v".`, service.Name))
	}
	return nil
}

func GetRemoveRules(config configuration.Configuration) RemoveRules {
	removeRules := make(RemoveRules)

	go func() {
		// wait for this channel to get any value
		<-removeRules

		logs.PrintInfo("Removing firewall rules...")
		for _, service := range services.Services {
			cmd := getIPTablesCommand(config, service, false)

			logs.PrintDebug(fmt.Sprintf(`Removing firewall rules for service "%v" with command "%v"...`, service.Name, cmd.String()))
			err := cmd.Run()
			if err != nil {
				logs.PrintError(fmt.Sprintf(`Failed to remove firewall rules for service "%v".`, service.Name))
			} else {
				logs.PrintInfo(fmt.Sprintf(`Removed firewall rules for service "%v".`, service.Name))
			}
		}

		os.Exit(0)
	}()

	return removeRules
}
