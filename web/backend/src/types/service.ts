import type { CerberoRegexes } from "./regex"

/**
 * @swagger
 * components:
 *  schemas:
 *    CerberoService:
 *      type: object
 *      properties:
 *        name:
 *          description: The name of the service
 *          type: string
 *        nfq:
 *          description: The nfq id of the service
 *          type: number
 *          minimum: 100
 *          maximum: 199
 *        port:
 *          description: The port number of the service
 *          type: number
 *          minimum: 1
 *          maximum: 65535
 *        protocol:
 *          description: The protocol used by the service
 *          type: string
 */

export type CerberoService = {
  name: string
  nfq: number
  port: number
  protocol: "tcp" | "udp"
  regexes?: CerberoRegexes[]
}

export type CerberoServiceCreate = Omit<CerberoService, "regexes"> & {
  regexes?: string[]
}

/**
 * @swagger
 * components:
 *   schemas:
 *    CerberoSetupResponse:
 *      type: object
 *      properties:
 *        isSetupDone:
 *          type: boolean
 */
