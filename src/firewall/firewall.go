package firewall

import (
	"cerbero3/firewall/firewallnfq"
	"cerbero3/firewall/rules"
	"cerbero3/services"
)

func Start(rr rules.RemoveRules) {
	for _, service := range services.Services {
		go firewallnfq.StartFirewallForService(rr, service)
	}
}
