import { DeleteObjectCommand } from '@aws-sdk/client-s3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()

  const body = await readBody(event)
  const { url } = body

  const fileKey = url.split('/').pop()

  if (!fileKey) {
    throw createError({
      statusCode: 400,
      statusMessage: 'กรุณาระบุชื่อไฟล์ (key) ที่ต้องการลบ',
    })
  }

  try {
    const s3 = useR2()

    await s3.send(new DeleteObjectCommand({
      Bucket: config.r2BucketName,
      Key: fileKey,
    }))


    return {
      success: 'success',
      message: `ลบภาพเรียบร้อยแล้ว`,
    }
  } catch (error: any) {
    throw createError({
      statusCode: 500,
      statusMessage: 'ลบไฟล์ไม่สำเร็จ: ' + error.message,
    })
  }
})