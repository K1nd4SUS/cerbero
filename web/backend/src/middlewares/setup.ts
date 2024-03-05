import { NextFunction, Request, Response } from "express"
import { Database } from "../database/db"

/**
 * Use this middleware on routes that require cerbero to have already been setup.
 */
export default async function setupMiddleware(
  req: Request,
  res: Response,
  next: NextFunction
) {
  const redis = Database.getInstance()
  const servicesKeys = await redis.keys("services:*")

  const isSetupDone = servicesKeys.length > 0

  if(isSetupDone) {
    next()
  }
  else {
    return res.status(404).json({
      error: "Cerbero has not been setup"
    })
  }
}
