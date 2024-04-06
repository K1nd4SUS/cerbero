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

  // The order of the regexes MATTERS!
  // If the regexes stay the same but their order changes the firewall is condidered NOT SYNCED!
  return res.json({
    isConnected: isFirewallConnected,
    isSynced: JSON.stringify(firewallConfig) === JSON.stringify(config)
  })
})

// Trigger a firewall configuration update
firewallRoute.post("/", async (req, res) => {
  const redis = Database.getInstance()

  if(!isFirewallConnected) {
    return res.status(409).json({
      error: "The firewall is not connected, can't update the configuration"
    })
  }

  // Update the "firewall config" in the db
  // Take the "web config" and override the firewall one
  const firewallKeys = await redis.keys("firewall:*")
  for(const firewallKey of firewallKeys) {
    await redis.del(firewallKey)
  }

  const servicesKeys = await redis.keys("services:*")
  const regexesKeys = await redis.keys("regexes:*:active")

  for(const serviceKey of servicesKeys) {
    await redis.copy(serviceKey, `firewall:${serviceKey}`)
  }

  for(const regexesKey of regexesKeys) {
    const regexesKeyParts = regexesKey.split(":")
    const newRegexesKey = regexesKeyParts.filter((_, i) => i !== regexesKeyParts.length - 1).join(":")

    await redis.copy(regexesKey, `firewall:${newRegexesKey}`)
  }


  cerberoEventEmitter.emit("cerberoConfigUpdate") // FIX: an update could fail inside of this and we would have updated the db anyways

  return res.status(204).end()
})

export default firewallRoute

