import { Router } from "express"
import swaggerJSDoc from "swagger-jsdoc"
import swaggerUi from "swagger-ui-express"

const swaggerRoute = Router()

/**
 * @swagger
 * /api/swagger:
 *  get:
 *    tags:
 *      - Swagger
 *    summary: Serves the Swagger UI
 *    description: Serves the Swagger UI
 *    responses:
 *      200:
 *        description: The UI was served successfully
 *        content:
 *          text/html:
 *            schema:
 */
swaggerRoute.get("/", swaggerUi.setup(swaggerJSDoc({
  definition: {
    openapi: "3.0.0",
    info: {
      title: "Cerbero API",
      version: "0.0.0"
    }
  },
  apis: [
    "./src/**/*.ts"
  ]
})))

export default swaggerRoute
