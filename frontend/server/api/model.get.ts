// server/api/model.get.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
  try {
    // รัน SQL Query โดยใช้ sql tagged template
    const models = await sql<{ modelname: string }[]>`
      SELECT * FROM models
    `
    return models
  } catch (error) {
    throw createError({
      statusCode: 500,
      statusMessage: 'Database connection error: ' + error.message,
    })
  }
})