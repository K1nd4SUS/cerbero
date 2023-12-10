import express from "express"

import "./config/env"
import "./utils/logs"

import logger from "./middlewares/logger"

const api = express()

api.use(logger)

api.get("/", (req, res) => {
  return res.json("The server is up and running.")
})

api.listen(process.env.API_PORT, () => {
  console.info(`API listening on port ${process.env.API_PORT}`)
})
