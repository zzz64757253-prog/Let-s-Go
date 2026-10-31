// https://nuxt.com/docs/api/configuration/nuxt-config

export default defineNuxtConfig({
  routeRules: {
    '/server/**': { proxy: 'http://localhost:8080/api/**' }
  },
  ssr: false,
  app: {
    head: {
      link: [
        { rel: 'icon', type: 'image/png', href: '/IconBS.png' }
      ]
    }
  },
  modules: ['@nuxtjs/tailwindcss', '@nuxt/icon/module'],
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  
  runtimeConfig: {
    // ค่าเหล่านี้จะถูกอ่านจาก .env อัตโนมัติและเข้าถึงได้เฉพาะฝั่ง Server เท่านั้น
    r2Endpoint: process.env.R2_ENDPOINT,
    r2AccessKeyId: process.env.R2_ACCESS_KEY_ID,
    r2SecretAccessKey: process.env.R2_SECRET_ACCESS_KEY,
    r2BucketName: process.env.R2_BUCKET_NAME,
    r2PublicDomain: process.env.R2_PUBLIC_DOMAIN,
    openaiApiKey: process.env.OPENAI_API_KEY,
  }
})
