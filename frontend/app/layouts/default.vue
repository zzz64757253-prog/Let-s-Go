<script setup lang="ts">
const isMobile = useState('isMobile', () => false)

onMounted(() => {
  const media = window.matchMedia('(max-width: 768px)')
  isMobile.value = media.matches

  // ดักจับเมื่อผู้ใช้ย่อ/ขยายหน้าจอ
  media.addEventListener('change', (e) => {
    isMobile.value = e.matches
  })
})
</script>

<template>
  <div class="min-h-screen w-full flex bg-gray-100" :class="isMobile ? 'flex-col' : 'flex-row'">
    <AppSidebar v-if="!isMobile"/>
    <AppNavbar v-else />
    
    <main class="w-full flex-1">
      <slot />
    </main>
  </div>
</template>
