import { Router } from "express"
import { Database } from "../database/db"

const setupRoute = Router()

/**
 * @swagger
 * /api/setup:
 *  get:
 *    tags:
 *      - Setup
 *    summary: Returns whether the setup has been done or not
 *    description: Returns whether the setup has been done or not. Cerbero is considered "set-up" if 1 or more services exist
 *    responses:
 *      200:
 *        description: Returns whether the setup has been done or not
 *        content:
 *          application/json:
 *            schema:
 *              $ref: '#/components/schemas/CerberoSetupResponse'
 */
setupRoute.get("/", async (req, res) => {
  const redis = Database.getInstance()

  const servicesKeys = await redis.keys("services:*")

  return res.send({
    isSetupDone: servicesKeys.length > 0
  })
})

export default setupRoute
