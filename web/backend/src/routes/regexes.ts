import { Router } from "express"
import { z } from "zod"
import { Database } from "../database/db"
import { CerberoServiceCreate } from "../types/service"

const regexesRoute = Router()

regexesRoute.get("/:nfq", async (req, res) => {
  const redis = Database.getInstance()
  const { nfq } = req.params

  if(!parseInt(nfq)) {
    return res.status(400).json({
      error: "The provided nfq id is not a number"
    })
  }

  // Check that the service exists before fetching regexes
  const serviceKeys = await redis.keys(`services:${nfq}`)

  if(!serviceKeys || serviceKeys.length <= 0) {
    return res.status(400).json({
      error: "You are trying to fetch regexes of a service that doesn't exist"
    })
  }

  const activeRegexes = await redis.sMembers(`regexes:${nfq}:active`)
  const inactiveRegexes = await redis.sMembers(`regexes:${nfq}:inactive`)

  return res.json({
    regexes: {
      active: activeRegexes,
      inactive: inactiveRegexes
    }
  })
})

regexesRoute.post("/:nfq", async (req, res) => {
  const redis = Database.getInstance()
  const { nfq } = req.params

  const bodySchema = z.object({
    regexes: z.array(z.string())
  })

  let typeValidatedBody: Required<Pick<CerberoServiceCreate, "regexes">>

  try {
    typeValidatedBody = bodySchema.parse(req.body)
  }
  catch(e) {
    return res.status(400).json({
      error: e
    })
  }

  if(!parseInt(nfq)) {
    return res.status(400).json({
      error: "The provided nfq id is not a number"
    })
  }

  // Check that the service exists before adding regexes
  const serviceKeys = await redis.keys(`services:${nfq}`)

  if(!serviceKeys || serviceKeys.length <= 0) {
    return res.status(400).json({
      error: "You are trying to add regexes to a service that doesn't exist"
    })
  }

  // Regexes are considered active by default
  await redis.sAdd(`regexes:${nfq}:active`, typeValidatedBody.regexes)

  const activeRegexes = await redis.sMembers(`regexes:${nfq}:active`)
  const inactiveRegexes = await redis.sMembers(`regexes:${nfq}:inactive`)

  return res.status(201).json({
    regexes: {
      active: activeRegexes,
      inactive: inactiveRegexes
    }
  })
})

export default regexesRoute
