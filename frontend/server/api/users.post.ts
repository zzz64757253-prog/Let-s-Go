// server/api/users.post.ts
import { sql } from '../utils/db'

export default defineEventHandler(async (event) => {
  try {
    // 1. ใช้ event อ่านข้อมูล JSON ที่ส่งมาจากหน้าบ้าน
    const body = await readBody(event)
    const { name, email } = body

    // 2. ตรวจสอบข้อมูลเบื้องต้น (Validation)
    if (!name || !email) {
      throw createError({
        statusCode: 400, // 400 = Bad Request (หน้าบ้านส่งข้อมูลมาไม่ครบ)
        statusMessage: 'กรุณากรอก name และ email ให้ครบถ้วน',
      })
    }

    // 3. บันทึกลง PostgreSQL ผ่านตัวแปร sql
    const newUser = await sql`
      INSERT INTO users (name, email)
      VALUES (${name}, ${email})
      RETURNING *
    `

    // 4. ส่งข้อมูลแถวใหม่ที่เพิ่มสำเร็จกลับไปให้หน้าบ้าน
    return {
      success: true,
      data: newUser[0],
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