import { Router } from "express"
import { cerberoEventEmitter, isFirewallConnected } from "../socket/socket"

const firewallRoute = Router()

firewallRoute.get("/", (req, res) => {
  // Get the last config that was sent from the db and
  // compare it to the current web config

  return res.json({
    isConnected: isFirewallConnected,
    isSynced: true
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

