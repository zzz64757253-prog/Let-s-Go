<script setup lang="ts">
const route = useRoute()

interface MenuType {
  menu_id: number
  name_menu: string
  menu_url: string | null
  parent_id: number | null
  menu_icon: string | null
}

interface MenuApiResponse {
  mainMenus: MenuType[]
  subMenus: MenuType[]
}

const { data: menusData } = await useFetch<MenuApiResponse>('/server/menu')

const isCollapsed = ref(true)
const isLogoHovered = ref(false)
const hoveredMenuId = ref<number | null>(null)
const openedMenuId = ref<number | null>(null)
const pinnedFlyoutId = ref<number | null>(null)

const defaultLogo = '/images/logo.png'
const hoverLogo = '/images/X.png'

// เปลี่ยนเมนู active ตาม URL โดยอัตโนมัติ (รวมหน้า submenu)
const activeSubmenu = computed(() =>
  menusData.value?.subMenus.find(item => item.menu_url === route.path) ?? null
)

const activeMainMenuId = computed<number | null>(() => {
  const mainMenu = menusData.value?.mainMenus.find(item => item.menu_url === route.path)
  return mainMenu?.menu_id ?? activeSubmenu.value?.parent_id ?? null
})

// เปิดหมวดที่กำลังใช้งานเมื่อเข้าหน้าหรือเปลี่ยนเส้นทาง
watch(activeMainMenuId, (menuId) => {
  if (menuId !== null) openedMenuId.value = menuId
}, { immediate: true })

// submenu ต้องเป็นรายการลูกของ menu นั้นเท่านั้น
function subMenusFor(menuId: number) {
  return (menusData.value?.subMenus ?? []).filter(item => item.parent_id === menuId && item.menu_url)
}

function toggleSubmenu(menuId: number) {
  if (isCollapsed.value) {
    pinnedFlyoutId.value = pinnedFlyoutId.value === menuId ? null : menuId
    return
  }
  openedMenuId.value = openedMenuId.value === menuId ? null : menuId
}

watch(() => route.path, () => {
  pinnedFlyoutId.value = null
})

function isActiveSubmenu(menuId: number) {
  return activeSubmenu.value?.menu_id === menuId
}
</script>

<template>
  <!-- แสดงเฉพาะจอคอมพิวเตอร์ (>= lg) -->
  <aside
    class="fixed inset-y-0 left-0 z-50 hidden flex-col border-r border-gray-200 bg-white shadow-sm transition-[width] duration-300 lg:flex"
    :class="isCollapsed ? 'w-16' : 'w-64'"
  >
    <!-- BRAND -->
    <div
      class="flex h-[76px] shrink-0 items-center border-b border-gray-100"
      :class="isCollapsed ? 'justify-center px-2' : 'justify-between px-4'"
    >
      <button
        type="button"
        class="flex h-11 w-11 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-blue-100 bg-blue-50 transition hover:bg-blue-100"
        :aria-label="isCollapsed ? 'ขยายเมนู' : 'ย่อเมนู'"
        :title="isCollapsed ? 'ขยายเมนู' : 'ย่อเมนู'"
        @mouseenter="isLogoHovered = true"
        @mouseleave="isLogoHovered = false"
        @click="isCollapsed = !isCollapsed; hoveredMenuId = null"
      >
        <img
          :src="isLogoHovered ? hoverLogo : defaultLogo"
          alt="โลโก้เว็บไซต์"
          class="h-full w-full object-contain"
        >
      </button>

      <div v-if="!isCollapsed" class="min-w-0 flex-1 px-3">
        <h1 class="truncate text-base font-bold text-gray-800">AkkaraponKi</h1>
        <p class="text-[11px] text-gray-400"></p>
      </div>

      <button
        v-if="!isCollapsed"
        type="button"
        aria-label="ย่อ Sidebar"
        class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-gray-400 hover:bg-gray-100 hover:text-gray-700"
        @click="isCollapsed = true"
      >
        <Icon name="heroicons:chevron-double-left-20-solid" class="h-4 w-4" />
      </button>
    </div>

    <!-- MENU -->
    <nav class="flex-1 space-y-1 py-4" aria-label="เมนูนำทางสำหรับคอมพิวเตอร์">
      <div
        v-for="menu in menusData?.mainMenus ?? []"
        :key="menu.menu_id"
        class="relative px-2"
        @mouseenter="hoveredMenuId = menu.menu_id"
        @mouseleave="hoveredMenuId = null; pinnedFlyoutId = null"
      >
        <!-- เมนูที่กดนำทางได้ทันที -->
        <NuxtLink
          v-if="menu.menu_url"
          :to="menu.menu_url"
          :title="isCollapsed ? menu.name_menu : undefined"
          :aria-current="route.path === menu.menu_url ? 'page' : undefined"
          class="relative flex h-11 items-center gap-3 rounded-lg text-sm transition-colors duration-200"
          :class="[
            isCollapsed ? 'justify-center px-1' : 'px-3',
            activeMainMenuId === menu.menu_id
              ? 'bg-blue-50 font-semibold text-blue-600'
              : 'text-gray-500 hover:bg-gray-50 hover:text-gray-800'
          ]"
        >
          <Icon :name="`heroicons:${menu.menu_icon || 'squares-2x2'}-20-solid`" class="h-5 w-5 shrink-0" />
          <span v-if="!isCollapsed" class="min-w-0 flex-1 truncate">{{ menu.name_menu }}</span>
          <span v-if="activeMainMenuId === menu.menu_id && isCollapsed" class="absolute right-0 h-5 w-[3px] rounded bg-blue-400"></span>
        </NuxtLink>

        <!-- เมนูแม่ที่มีเมนูย่อย -->
        <button
          v-else
          type="button"
          :title="isCollapsed ? menu.name_menu : undefined"
          :aria-expanded="isCollapsed ? hoveredMenuId === menu.menu_id || pinnedFlyoutId === menu.menu_id : openedMenuId === menu.menu_id"
          class="relative flex h-11 w-full items-center gap-3 rounded-lg text-left text-sm transition-colors duration-200"
          :class="[
            isCollapsed ? 'justify-center px-1' : 'px-3',
            activeMainMenuId === menu.menu_id
              ? 'bg-blue-50 font-semibold text-blue-600'
              : 'text-gray-500 hover:bg-gray-50 hover:text-gray-800'
          ]"
          @click="toggleSubmenu(menu.menu_id)"
        >
          <Icon :name="`heroicons:${menu.menu_icon || 'folder'}-20-solid`" class="h-5 w-5 shrink-0" />
          <span v-if="!isCollapsed" class="min-w-0 flex-1 truncate">{{ menu.name_menu }}</span>
          <Icon
            v-if="!isCollapsed"
            name="heroicons:chevron-down-20-solid"
            class="h-4 w-4 shrink-0 transition-transform duration-200"
            :class="openedMenuId === menu.menu_id ? 'rotate-180' : ''"
          />
          <span v-if="activeMainMenuId === menu.menu_id && isCollapsed" class="absolute right-0 h-5 w-[3px] rounded bg-blue-400"></span>
        </button>

        <!-- โหมดขยาย: submenu อยู่ด้านล่างของเมนูแม่ -->
        <div
          v-if="!isCollapsed && !menu.menu_url && openedMenuId === menu.menu_id"
          class="ml-5 mt-1 space-y-1 border-l border-gray-200 pl-3"
        >
          <NuxtLink
            v-for="subMenu in subMenusFor(menu.menu_id)"
            :key="subMenu.menu_id"
            :to="subMenu.menu_url || '/'"
            class="flex items-center gap-2 rounded-md px-3 py-2 text-xs transition-colors"
            :class="isActiveSubmenu(subMenu.menu_id)
              ? 'bg-blue-50 font-semibold text-blue-600'
              : 'text-gray-500 hover:bg-gray-50 hover:text-gray-800'"
          >
            <Icon :name="`heroicons:${subMenu.menu_icon || 'document'}-20-solid`" class="h-4 w-4 shrink-0" />
            <span class="truncate">{{ subMenu.name_menu }}</span>
          </NuxtLink>
        </div>

        <!-- โหมดย่อ: submenu ลอยด้านข้างเมื่อ hover หรือคลิก -->
        <div
          v-if="isCollapsed && !menu.menu_url && (hoveredMenuId === menu.menu_id || pinnedFlyoutId === menu.menu_id)"
          class="absolute left-full top-0 z-[60] w-56 rounded-r-lg rounded-bl-lg border border-gray-200 bg-white p-2 shadow-lg"
        >
          <p class="border-b border-gray-100 px-3 py-2 text-xs font-semibold text-gray-800">
            {{ menu.name_menu }}
          </p>
          <NuxtLink
            v-for="subMenu in subMenusFor(menu.menu_id)"
            :key="subMenu.menu_id"
            :to="subMenu.menu_url || '/'"
            class="mt-1 flex items-center gap-2 rounded-md px-3 py-2.5 text-sm transition-colors"
            :class="isActiveSubmenu(subMenu.menu_id)
              ? 'bg-blue-50 font-semibold text-blue-600'
              : 'text-gray-500 hover:bg-gray-50 hover:text-gray-800'"
          >
            <Icon :name="`heroicons:${subMenu.menu_icon || 'document'}-20-solid`" class="h-4 w-4 shrink-0" />
            <span class="truncate">{{ subMenu.name_menu }}</span>
          </NuxtLink>
        </div>
      </div>
    </nav>

    <div class="shrink-0 border-t border-gray-100 px-3 py-4 text-center">
      <span v-if="!isCollapsed" class="text-[11px] text-gray-400">© 2026 Akkarapon</span>
      <Icon v-else name="heroicons:ellipsis-horizontal-20-solid" class="h-5 w-5 text-gray-300" />
    </div>
  </aside>

  <!-- Spacer: ใช้งานได้เมื่อ layout หลักจัดวางแบบ flex -->
  <div
    class="hidden shrink-0 transition-[width] duration-300 lg:block"
    :class="isCollapsed ? 'w-16' : 'w-64'"
    aria-hidden="true"
  ></div>
</template>
