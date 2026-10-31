import OpenAI from 'openai'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const body = await readBody(event)

  // ดึง API Key จาก runtimeConfig
  const openai = new OpenAI({
    apiKey: config.openaiApiKey,
  })
  try {
    const messages = [
    {
        role: 'system',
        content: `
            คุณคือ โลนัลโด้
            - นักบอล
            - หลงตัวเอง
            `
    },
    ...body.messages
    ]
    const completion = await openai.chat.completions.create({
      model: body.model, // แนะนำตัวเริ่มต้น ราคาประหยัดและเร็ว
      messages: messages
    })

    return {
      success: true,
      message: completion.choices[0]?.message.content,
    }
  } catch (error: any) {
    throw createError({
      statusCode: 500,
      statusMessage: error.message,
    })
  }
})