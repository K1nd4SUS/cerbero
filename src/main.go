package main

import (
	"cerbero3/configuration"
	"cerbero3/credentials"
	"cerbero3/firewall"
	"cerbero3/firewall/rules"
	"cerbero3/interrupt"
	"cerbero3/logs"
	"cerbero3/services"
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

	// TODO: handle metrics

	logs.PrintInfo("Loading user flags...")
	config := configuration.GetFlagsConfiguration()
	if err = configuration.CheckValues(config); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	logs.Configure(config)
	// debug logs can be displayed only after this commented line
	logs.PrintInfo("Loaded user flags.")

	logs.PrintInfo("Loading services configuration...")
	if err = services.Load(config.ConfigurationFile); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	if err = services.CheckServicesValues(); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	logs.PrintInfo("Loaded services configuration.")

	logs.PrintInfo("Setting firewall rules...")
	if err = rules.SetRules(config); err != nil {
		logs.PrintCritical(err.Error())
		os.Exit(1)
	}
	logs.PrintInfo("Set firewall rules.")

	// whenever "rr <- true" is called, the program is stopped
	rr := rules.GetRemoveRules(config)

	logs.PrintInfo("Starting firewall...")
	firewall.Start(rr)
	logs.PrintInfo("Started firewall.")

	interrupt.Listen(rr)

	// TODO: create alert for when the config file is edited

	// hang this thread; the program will be closed by other means,
	// being an os interrupt signal or a crash
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
