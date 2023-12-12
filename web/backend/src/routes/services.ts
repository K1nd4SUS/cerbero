import { Router } from "express"
import { z } from "zod"
import { Database } from "../database/db"
import type { CerberoServiceCreate } from "../types/service"

const servicesRoute = Router()

// TODO: convert nfq and port to numbers (since redis returns everything as a string)

servicesRoute.get("/", async (req, res) => {
  const redis = Database.getInstance()

  const servicesKeys = await redis.keys("services:*")
  const services = []

  for(const serviceKey of servicesKeys) {
    const service = await redis.hGetAll(serviceKey)

    services.push(service)
  }

  return res.json(services)
})

servicesRoute.post("/", async (req, res) => {
  const redis = Database.getInstance()

  const bodySchema = z.object({
    name: z.string(),
    nfq: z.number(),
    port: z.number(),
    protocol: z.literal("tcp").or(z.literal("udp")),
    regexes: z.array(z.string()).optional()
  })

  let typeValidatedBody: CerberoServiceCreate

  try {
    typeValidatedBody = bodySchema.parse(req.body)
  }
  catch(e) {
    return res.status(400).json({
      error: e
    })
  }

  // The service was created with default regexes
  if(typeValidatedBody.regexes && typeValidatedBody.regexes.length > 0) {
    await redis.sAdd(`regexes:${typeValidatedBody.nfq}:active`, typeValidatedBody.regexes)
  }

  const newService = {
    name: typeValidatedBody.name,
    nfq: typeValidatedBody.nfq,
    port: typeValidatedBody.port,
    protocol: typeValidatedBody.protocol
  }

  await redis.hSet(`services:${typeValidatedBody.nfq}`, newService)

  return res.status(201).json({
    ...newService,
    regexes: {
      active: typeValidatedBody.regexes,
      inactive: []
    }
  })
})

servicesRoute.get("/:nfq", async (req, res) => {
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

  return res.json({
    ...service,
    regexes: {
      active: activeRegexes,
      inactive: inactiveRegexes
    }
  })
})

export default servicesRoute
