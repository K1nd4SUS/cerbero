import express from "express"

import "./config/env"
import "./utils/logs"

import { Database } from "./database/db"
import logger from "./middlewares/logger"
import regexesRoute from "./routes/regexes"
import servicesRoute from "./routes/services"
import setupRoute from "./routes/setup"

// Instantiate redis connection
Database
  .getInstance()
  .connect()
  .then(() => console.info(`Connected to db ${process.env.REDIS_URL}`))
  .catch(() => {
    console.error(`The API couldn't connect to the db ${process.env.REDIS_URL}`)
    process.exit(1)
  })

const api = express()

api.use(express.json())
api.use(logger)

api.use("/api/regexes", regexesRoute)
api.use("/api/services", servicesRoute)
api.use("/api/setup", setupRoute)

api.listen(process.env.API_PORT, () => {
  console.info(`API listening on port ${process.env.API_PORT}`)
})
