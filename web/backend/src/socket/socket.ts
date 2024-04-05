import EventEmitter from "events"
import { createServer } from "net"
import { Database } from "../database/db"

export let isFirewallConnected = false
export const cerberoEventEmitter = new EventEmitter()

/**
 * Fetches the services and regexes from the db
 * and returns the json config to feed to the firewall
 */
export async function buildConfiguration() {
  const redis = Database.getInstance()

  const config = []
  const servicesKeys = await redis.keys("services:*")

  for(const serviceKey of servicesKeys) {
    const nfq = serviceKey.split(":")[1]
    const service = await redis.hGetAll(serviceKey)
    const regexes = await redis.sMembers(`regexes:${nfq}:active`)

    config.push({
      chain: "OUTPUT", // TODO: set from db
      name: service.name,
      nfq: parseInt(service.nfq),
      port: parseInt(service.port),
      protocol: service.protocol as "tcp" | "udp",
      regexes: regexes
    })
  }

  return config
}

/**
 * Only one connection at a time can be established
 * Once the first client is connected (firewall) the server stops listening for new connections
 * In case the connection is closed, the server starts listening again
 */
const socketServer = createServer()

socketServer.on("listening", () => {
  console.info(`Socket server listening on port ${process.env.SOCKET_PORT}`)
})

socketServer.on("connection", async socket => {
  const { remoteAddress, remotePort } = socket

  console.info(`A connection with ${remoteAddress}:${remotePort} was established`)

  isFirewallConnected = true
  socketServer.close()

  console.info("The socket server has stopped listening for connections")

  socket.once("data", data => {
    console.info(`Initialization string: ${data.toString()}`)
  })

  cerberoEventEmitter.on("cerberoConfigUpdate", async () => {
    if(!isFirewallConnected) {
      console.error("Tried to update the config but the firewall is disconnected")
      return
    }

    const newConfig = await buildConfiguration()
    const encodedConfig = btoa(JSON.stringify(newConfig))

    socket.write(encodedConfig + "\n")
  })

  cerberoEventEmitter.emit("cerberoConfigUpdate")

  socket.on("close", hadError => {
    isFirewallConnected = false

    if(hadError) {
      console.error(`The connection with ${remoteAddress}:${remotePort} had an error`)
    }

    console.info(`The connection with ${remoteAddress}:${remotePort} was closed`)

    socketServer.listen(process.env.SOCKET_PORT)
  })
})

export default socketServer

