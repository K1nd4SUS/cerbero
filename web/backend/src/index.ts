import express from "express"
import "./config/env.ts"

const server = express()

server.get("/", (req, res) => {
  console.log(req)

  return res.send("The server is up and running.")
})

server.listen(process.env.SERVER_PORT, () => {
  console.log(`Server listening on port ${process.env.SERVER_PORT}`)
})
