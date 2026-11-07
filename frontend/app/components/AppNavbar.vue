<script setup lang="ts">
const route = useRoute()

interface MenuType {
  menu_id: number
  name_menu: string
  menu_url: string
  menu_icon: string
}

interface MenuResponse {
  menus: MenuType[]
}

const { data: menus } = await useFetch<MenuResponse>('/server/mobilemenu')

const isActive = (url: string) => route.path === url
</script>

<template>
  <!-- ใช้บนมือถือและแท็บเล็ต (< lg) -->
  <nav
    v-if="menus?.menus?.length"
    aria-label="เมนูนำทางสำหรับมือถือ"
    class="fixed inset-x-0 top-0 z-50 h-16 border-b border-gray-200 bg-white/95 shadow-sm backdrop-blur-md lg:hidden"
  >
    <div class="mobile-menu-scroll flex h-full items-stretch gap-1 overflow-x-auto px-2 sm:px-4">
      <NuxtLink
        v-for="menu in menus.menus"
        :key="menu.menu_id"
        :to="menu.menu_url"
        :title="menu.name_menu"
        :aria-label="menu.name_menu"
        :aria-current="isActive(menu.menu_url) ? 'page' : undefined"
        class="relative flex min-w-16 flex-1 flex-col items-center justify-center gap-1 rounded-md px-2 pt-1 text-gray-400 transition-colors duration-200 hover:bg-blue-50 hover:text-blue-500"
        :class="isActive(menu.menu_url) ? 'text-blue-500' : ''"
      >
        <Icon
          :name="`heroicons:${menu.menu_icon || 'squares-2x2'}-20-solid`"
          class="h-5 w-5 shrink-0"
        />
        <span
          class="max-w-full truncate text-[10px] leading-4"
          :class="isActive(menu.menu_url) ? 'font-semibold' : 'font-medium'"
        >
          {{ menu.name_menu }}
        </span>
        <span
          class="absolute bottom-0 left-1/2 h-[3px] rounded-full bg-blue-400 transition-all duration-300 -translate-x-1/2"
          :class="isActive(menu.menu_url) ? 'w-8 opacity-100' : 'w-0 opacity-0'"
        ></span>
      </NuxtLink>
    </div>
  </nav>

  <!-- เว้นที่ให้ navbar ด้านบนไม่ทับเนื้อหา -->
  <div class="h-16 lg:hidden" aria-hidden="true"></div>
</template>

<style scoped>
.mobile-menu-scroll {
  scrollbar-width: none;
}
.mobile-menu-scroll::-webkit-scrollbar {
  display: none;
}
</style>
