import { Router } from "express"
import { z } from "zod"
import { Database } from "../database/db"
import setupMiddleware from "../middlewares/setup"

const regexesRoute = Router()

/**
 * @swagger
 * /api/regexes/{nfq}:
 *  get:
 *    tags:
 *      - Regexes
 *    summary: Get all the regexes of a service (by its nfq id)
 *    description: Get all the regexes of a service (by its nfq id)
 *    parameters:
 *      - $ref: '#/components/parameters/nfq'
 *    responses:
 *      200:
 *        description: The regexes were returned successfully
 *        content:
 *          application/json:
 *            schema:
 *              $ref: '#/components/schemas/CerberoRegexes'
 *      400:
 *        $ref: '#/components/responses/BadRequest'
 */
regexesRoute.get("/:nfq", setupMiddleware, async (req, res) => {
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

/**
 * @swagger
 * /api/regexes/{nfq}:
 *  post:
 *    tags:
 *      - Regexes
 *    summary: Add regexes to a service (by its nfq id)
 *    description: Add regexes to a service (by its nfq id)
 *    parameters:
 *      - $ref: '#/components/parameters/nfq'
 *    requestBody:
 *      required: true
 *      content:
 *        application/json:
 *          schema:
 *            $ref: '#/components/schemas/Regexes'
 *    responses:
 *      201:
 *        description: The regexes were added successfully
 *        content:
 *          application/json:
 *            schema:
 *              $ref: '#/components/schemas/CerberoRegexes'
 *      400:
 *        $ref: '#/components/responses/BadRequest'
 */
regexesRoute.post("/:nfq", setupMiddleware, async (req, res) => {
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

/**
 * @swagger
 * /api/regexes/{nfq}:
 *  put:
 *    tags:
 *      - Regexes
 *    summary: Modify the state of a regex (active/inactive)
 *    description: Modify the state of a regex (active/inactive)
 *    parameters:
 *      - $ref: '#/components/parameters/nfq'
 *      - $ref: '#/components/parameters/reghex'
 *    requestBody:
 *      required: true
 *      content:
 *        application/json:
 *          schema:
 *            $ref: '#/components/schemas/PutRegex'
 *    responses:
 *      200:
 *        description: The regex was modified successfully
 *        content:
 *          application/json:
 *            schema:
 *              $ref: '#/components/schemas/PutRegex'
 */
regexesRoute.put("/:nfq", setupMiddleware, async (req, res) => {
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

/**
 * @swagger
 * /api/regexes/{nfq}:
 *  delete:
 *    tags:
 *      - Regexes
 *    summary: Delete a regex
 *    description: Delete a regex
 *    parameters:
 *      - $ref: '#/components/parameters/nfq'
 *    responses:
 *      204:
 *        description: The regex was deleted successfully
 *      400:
 *        $ref: '#/components/responses/BadRequest'
 */
regexesRoute.delete("/:nfq", setupMiddleware, async (req, res) => {
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
