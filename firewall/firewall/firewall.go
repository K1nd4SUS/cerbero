package firewall

import (
	"cerbero/firewall/firewallnfq"
	"cerbero/firewall/rules"
	"cerbero/services"
)

func Start(rr rules.RemoveRules) {
	for index := range services.Services {
		// we use the index so that the services can be
		// updated at runtime, without worrying about
		// sending signals to each separate thread
		// that was previously started
		go firewallnfq.StartFirewallForService(rr, index)
	}
}
