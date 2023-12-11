import { createClient, RedisClientType } from "redis"

/*
 * Singleton for redis db connection
 * Access the redis client with: Database.getInstance()
 *
 * Example:
 * const redis = Database()
 *  .getInstance()
 *  .connect()
 */
export class Database {
  private static db: RedisClientType | undefined

  private constructor() {}

  public static getInstance() {
    if(!this.db) {
      this.db = createClient({
        url: process.env.REDIS_URL,
        database: 0
      })
    }

    return this.db
  }
}
