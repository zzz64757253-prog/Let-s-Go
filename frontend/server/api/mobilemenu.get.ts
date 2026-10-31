// server/api/content.get.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
    try {
        // รัน SQL Query โดยใช้ sql tagged template
        const menus = await sql`
            SELECT * FROM menu WHERE menu_url IS NOT NULL ORDER BY menu_id ASC
        `
        return menus
    } catch (error: any) {
        throw createError({
            statusCode: 500,
            statusMessage: 'Database connection error: ' + error.message,
        })
    }
})