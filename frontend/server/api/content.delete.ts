// server/api/content.delete.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
    try {
    // 1. ใช้ event อ่านข้อมูล JSON ที่ส่งมาจากหน้าบ้าน
        const body = await readBody(event)
        const { post_uid, pin } = body


        // 3. บันทึกลง PostgreSQL ผ่านตัวแปร sql
        const deletedPost = await sql`
        DELETE FROM posts WHERE post_uid = ${post_uid} AND pin = ${pin}
        RETURNING *
        `

        // 4. ส่งข้อมูลแถวใหม่ที่เพิ่มสำเร็จกลับไปให้หน้าบ้าน
        if (deletedPost.length === 0) {
            return {
                success: "failure",
                message: 'PIN ไม่ถูกต้อง!!',
            }
        }
        return {
            success: "success",
            message: 'ลบโพสต์สำเร็จแล้ว',
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