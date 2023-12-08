package main

import (
	"cerbero3/configuration"
	"cerbero3/credentials"
	"cerbero3/firewall"
	"cerbero3/firewall/rules"
	"cerbero3/interrupt"
	"cerbero3/logs"
	"cerbero3/metrics"
	"cerbero3/services"
	"errors"
	"fmt"
	"os"
	"sync"
)

func main() {
	isRoot, err := credentials.IsUserRoot()
	if err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	if !isRoot {
		logs.PrintCritical("This firewall must be run as root.")
		os.Exit(1)
	}

	logs.PrintInfo("Loading user flags...")
	config := configuration.GetFlagsConfiguration()
	if err = configuration.CheckValues(&config); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	logs.Configure(config)
	// debug logs can be displayed only after this commented line
	logs.PrintInfo("Loaded user flags.")

	logs.PrintInfo("Loading services configuration for the first time...")
	if configuration.IsConfigFileSet(config) {
		err = services.LoadConfigFile(config.ConfigurationFile)
	} else if configuration.IsCerberoSocketSet(config) {
		err = services.LoadCerberoSocket(config.CerberoSocketIP, config.CerberoSocketPort, 0)
	} else {
		err = errors.New("Neither a configuration file nor a socket were found.")
	}
	if err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	if err = services.CheckServicesValues(); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	logs.PrintInfo("Loaded services configuration for the first time.")

	logs.PrintInfo("Setting firewall rules...")
	if err = rules.SetRules(config); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	logs.PrintInfo("Set firewall rules.")

	logs.PrintInfo("Starting metrics server...")
	go metrics.StartServer(config)
	logs.PrintInfo(fmt.Sprintf("Started metrics server on port %v.", config.MetricsPort))

	// whenever "rr <- true" is called, the program is stopped
	rr := rules.GetRemoveRules(config)

	logs.PrintInfo("Starting firewall...")
	firewall.Start(rr)
	logs.PrintInfo("Started firewall.")

	interrupt.Listen(rr)

	// hang this thread; the program will be closed by other means,
	// being an os interrupt signal or a crash
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
