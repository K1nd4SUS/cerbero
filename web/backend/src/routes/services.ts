import { Router } from "express"
import { Database } from "../database/db"

const servicesRoute = Router()

servicesRoute.get("/", async (req, res) => {
  const redis = Database.getInstance()

  const servicesKeys = await redis.keys("services:*")
  const services = []

  for(const serviceKey of servicesKeys) {
    const service = await redis.hGetAll(serviceKey)

    services.push({ ...service, name: serviceKey.split(":")[1] })
  }

  return res.json(services)
})

export default servicesRoute
