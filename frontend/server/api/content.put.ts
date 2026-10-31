import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
    try {
        // 1. ใช้ event อ่านข้อมูล JSON ที่ส่งมาจากหน้าบ้าน
        const body = await readBody(event)
        const { post_uid, username, title, content, image_url, pin } = body


        const pin1: [{ pin: string }] = await sql`
        SELECT pin FROM posts 
        WHERE post_uid = ${post_uid}`
        if (pin !== pin1[0].pin) {
            return {
                success: 'failure',
                message: 'PIN ไม่ถูกต้อง',
            }
        }

        await sql`
        UPDATE posts SET username = ${username}, title = ${title}, content = ${content}, image_url = ${image_url}
        WHERE post_uid = ${post_uid}
        `
        return {
            success: 'success',
            message: 'แก้ไขโพสต์สำเร็จ',
        }

    } catch (error: any) {
        if (error.statusCode) throw error

        throw createError({
            statusCode: 500,
            statusMessage: 'Database Error: ' + error.message,
        })
    }
})