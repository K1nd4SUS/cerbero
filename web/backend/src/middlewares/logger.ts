import { Request, Response, NextFunction } from "express"

export default function logger(req: Request, res: Response, next: NextFunction) {
  console.log(
    req.ip,
    req.method,
    req.path,
    req.protocol,
    req.httpVersion
  )

  next()
}
