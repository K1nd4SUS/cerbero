import dotenv from "dotenv"
import { z } from "zod"

dotenv.config()

const env = z.object({
  API_PORT: z.string().transform(v => parseInt(v)),
  REDIS_URL: z.string(),
  SOCKET_PORT: z.string()
})

env.parse(process.env)

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace NodeJS {
    interface ProcessEnv extends z.infer<typeof env> {}
  }
}
