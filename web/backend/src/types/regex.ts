/**
 * @swagger
 * components:
 *  schemas:
 *    Regexes:
 *      type: object
 *      required:
 *        - regexes
 *      properties:
 *        regexes:
 *          type: array
 *          items:
 *            type: string
 *          description: The list of regexes to add
 *    CerberoRegexes:
 *      type: object
 *      properties:
 *        regexes:
 *          type: object
 *          properties:
 *            active:
 *              type: array
 *              items:
 *                type: string
 *              description: The active regexes
 *            inactive:
 *              type: array
 *              items:
 *                type: string
 *              description: The inactive regexes
 *          description: The list of regexes to add
 *    PutRegex:
 *      type: object
 *      properties:
 *        regex:
 *          type: string
 *        state:
 *          type: string
 */

export type CerberoRegexes = {
  active: string[]
  inactive: string[]
}
