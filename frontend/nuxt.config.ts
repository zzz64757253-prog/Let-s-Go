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
  
})
