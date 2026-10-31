// server/api/content.get.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
  try {
    // รัน SQL Query โดยใช้ sql tagged template
    const posts = await sql`
      SELECT * FROM posts ORDER BY post_date DESC
    `
    return posts
  } catch (error: any) {
    throw createError({
      statusCode: 500,
      statusMessage: 'Database connection error: ' + error.message,
    })
  }
})