package main

import (
	"cerbero3/configuration"
	"cerbero3/credentials"
	"cerbero3/logs"
	"cerbero3/services"
	"os"
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

	// logs.PrintWarning("This is a warning.")
	// logs.PrintError("This is an error.")
	// logs.PrintCritical("This is a critical error.")
	// logs.PrintInfo("This is an info.")
}
