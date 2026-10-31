// server/api/users.get.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
  try {
    // รัน SQL Query โดยใช้ sql tagged template
    const users = await sql`
      SELECT * FROM users LIMIT 10
    `
    return users
  } catch (error : any) {
    throw createError({
      statusCode: 500,
      statusMessage: 'Database connection error: ' + error.message,
    })
  }
})