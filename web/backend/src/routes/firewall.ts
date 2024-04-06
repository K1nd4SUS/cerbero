import { Router } from "express"
import { Database } from "../database/db"
import { buildConfiguration, cerberoEventEmitter, isFirewallConnected } from "../socket/socket"

const firewallRoute = Router()

firewallRoute.get("/", async (req, res) => {
  const redis = Database.getInstance()

  const config = await buildConfiguration()
  const firewallConfig = []
  const firewallServicesKeys = await redis.keys("firewall:services:*")

  for(const firewallServiceKey of firewallServicesKeys) {
    const nfq = firewallServiceKey.split(":")[2]
    const firewallService = await redis.hGetAll(firewallServiceKey)
    const regexes = await redis.sMembers(`firewall:regexes:${nfq}`)

    const parsedFirewallService = {
      chain: "OUTPUT",
      name: firewallService.name,
      nfq: parseInt(firewallService.nfq),
      port: parseInt(firewallService.port),
      protocol: firewallService.protocol as "tcp" | "udp",
      regexes: regexes
    }

    firewallConfig.push(parsedFirewallService)
  }

  const sortedConfig = config.sort((service1, service2) => {
    if(service1.name > service2.name) {
      return -1
    }

    return 1
  })

  const sortedFirewallConfig = firewallConfig.sort((service1, service2) => {
    if(service1.name > service2.name) {
      return -1
    }

    return 1
  })

  // The order of the regexes MATTERS!
  return res.json({
    isConnected: isFirewallConnected,
    isSynced: JSON.stringify(sortedFirewallConfig) === JSON.stringify(sortedConfig)
  })
})

// Trigger a firewall configuration update
firewallRoute.post("/", (req, res) => {
  if(!isFirewallConnected) {
    return res.status(409).json({
      error: "The firewall is not connected, can't update the configuration"
    })
  }

  cerberoEventEmitter.emit("cerberoConfigUpdate")

  return res.status(204).end()
})

export default firewallRoute

