import { Router } from "express"
import { z } from "zod"
import { Database } from "../database/db"

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

  let typeValidatedBody

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

regexesRoute.put("/:nfq", async (req, res) => {
  const redis = Database.getInstance()

  const { nfq } = req.params

  const querySchema = z.object({
    reghex: z.string(),
    state: z.literal("active").or(z.literal("inactive"))
  })

  const bodySchema = z.object({
    regex: z.string(),
    state: z.literal("active").or(z.literal("inactive"))
  })

  let typeValidatedQuery
  let typeValidatedBody

  try {
    typeValidatedQuery = querySchema.parse(req.query)
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

  const parsedRegex = Buffer
    .from(typeValidatedQuery.reghex, "hex")
    .toString("utf-8")
    .trim()

  const doesRegexExist = await redis.sIsMember(
    `regexes:${nfq}:${typeValidatedQuery.state}`,
    parsedRegex
  )

  if(!doesRegexExist) {
    return res.status(400).json({
      error: "The regex you are trying to edit doesn't exist"
    })
  }

  await redis.sRem(`regexes:${nfq}:${typeValidatedQuery.state}`, parsedRegex)
  await redis.sAdd(`regexes:${nfq}:${typeValidatedBody.state}`, typeValidatedBody.regex)

  return res.json(typeValidatedBody)
})

regexesRoute.delete("/:nfq", async (req, res) => {
  const redis = Database.getInstance()
  const { nfq } = req.params

  const querySchema = z.object({
    reghex: z.string(),
    state: z.literal("active").or(z.literal("inactive"))
  })

  let typeValidatedQuery

  try {
    typeValidatedQuery = querySchema.parse(req.query)
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

  const parsedRegex = Buffer
    .from(typeValidatedQuery.reghex, "hex")
    .toString("utf-8")
    .trim()

  const doesRegexExist = await redis.sIsMember(
    `regexes:${nfq}:${typeValidatedQuery.state}`,
    parsedRegex
  )

  if(!doesRegexExist) {
    return res.status(400).json({
      error: "The regex you are trying to delete doesn't exist"
    })
  }

  await redis.sRem(`regexes:${nfq}:${typeValidatedQuery.state}`, parsedRegex)

  return res.status(204).end()
})

export default regexesRoute
