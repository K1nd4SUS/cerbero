import { Router } from "express"
import { z } from "zod"
import { Database } from "../database/db"
import setupMiddleware from "../middlewares/setup"
import type { CerberoService } from "../types/service"

const servicesRoute = Router()

servicesRoute.get("/", setupMiddleware, async (req, res) => {
  const redis = Database.getInstance()

  const servicesKeys = await redis.keys("services:*")
  const services = []

  for(const serviceKey of servicesKeys) {
    const service = await redis.hGetAll(serviceKey)

    const parsedService: CerberoService = {
      name: service.name,
      nfq: parseInt(service.nfq),
      port: parseInt(service.port),
      protocol: service.protocol as "tcp" | "udp"
    }

    services.push(parsedService)
  }

  return res.json(services)
})

servicesRoute.post("/", async (req, res) => {
  const redis = Database.getInstance()

  // Check if cerbero has already been setup
  const servicesKeys = await redis.keys("services:*")

  if(servicesKeys.length > 0) {
    return res.status(400).json({
      error: "Cerbero has already been setup"
    })
  }

  const bodySchema = z.array(z.object({
    name: z.string(),
    nfq: z.number(),
    port: z.number(),
    protocol: z.literal("tcp").or(z.literal("udp")),
    regexes: z.array(z.string()).optional()
  }))

  let typeValidatedBody

  try {
    typeValidatedBody = bodySchema.parse(req.body)
  }
  catch(e) {
    return res.status(400).json({
      error: e
    })
  }

  for(const service of typeValidatedBody) {
    // The service was created with default regexes
    if(service.regexes && service.regexes.length > 0) {
      await redis.sAdd(`regexes:${service.nfq}:active`, service.regexes)
    }

    const newService = {
      name: service.name,
      nfq: service.nfq,
      port: service.port,
      protocol: service.protocol
    }

    await redis.hSet(`services:${service.nfq}`, newService)
  }

  return res.status(201).json({
    isSetupDone: true
  })
})

servicesRoute.get("/:nfq", setupMiddleware, async (req, res) => {
  const redis = Database.getInstance()
  const { nfq } = req.params

  if(!parseInt(nfq)) {
    return res.status(400).json({
      error: "The provided nfq id is not a number"
    })
  }

  const service = await redis.hGetAll(`services:${nfq}`)
  const activeRegexes = await redis.sMembers(`regexes:${nfq}:active`)
  const inactiveRegexes = await redis.sMembers(`regexes:${nfq}:inactive`)

  if(Object.entries(service).length === 0) {
    return res.status(404).json({
      error: "Service not found"
    })
  }

  const parsedService: CerberoService = {
    name: service.name,
    nfq: parseInt(service.nfq),
    port: parseInt(service.port),
    protocol: service.protocol as "tcp" | "udp"
  }

  return res.json({
    ...parsedService,
    regexes: {
      active: activeRegexes,
      inactive: inactiveRegexes
    }
  })
})

export default servicesRoute
