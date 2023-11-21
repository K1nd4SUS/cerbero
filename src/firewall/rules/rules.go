package rules

import (
	"cerbero3/configuration"
	"cerbero3/logs"
	"cerbero3/services"
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
		}()), config.Chain,
		"-p", service.Protocol,
		"--dport", fmt.Sprintf("%v", service.Port),
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
				logs.PrintError(err.Error())
			} else {
				logs.PrintInfo(fmt.Sprintf(`Removed firewall rules for service "%v".`, service.Name))
			}
		}

		os.Exit(0)
	}()

	return removeRules
}
