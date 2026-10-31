// server/api/users.post.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
    try {
    // 1. ใช้ event อ่านข้อมูล JSON ที่ส่งมาจากหน้าบ้าน
        const body = await readBody(event)
        const { username, title, content, image_url, pin } = body


        // 3. บันทึกลง PostgreSQL ผ่านตัวแปร sql
        await sql`
        INSERT INTO posts (username, title, content, image_url, pin)
        VALUES (${username}, ${title}, ${content}, ${image_url}, ${pin})
        `

        // 4. ส่งข้อมูลแถวใหม่ที่เพิ่มสำเร็จกลับไปให้หน้าบ้าน
        return {
            success: 'success',
            message: 'โพสต์สำเร็จ',
        }

    } catch (error: any) {
    // ถ้าเป็น Error ที่เราสร้างไว้ข้างบน ให้พ่นออกไปเลย
    if (error.statusCode) throw error

    // ถ้าเป็น Database error อื่นๆ ให้ส่ง 500
    throw createError({
        statusCode: 500,
        statusMessage: 'Database Error: ' + error.message,
    })
    }
})