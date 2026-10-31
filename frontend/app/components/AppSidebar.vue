<script setup lang="ts">
import { NuxtLink } from '#components'
const route = useRoute()
console.log(route.path) 
const hoveredMenuId = ref<number | null>(null)
const LogoHover = ref<boolean>(false)
const LogoActive = ref<boolean>(false)
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
const menu = ref<any>(null)
const { data: menusData} = await useFetch<MenuApiResponse>('/server/menu')

const ActiveMainmenu = ref<number | null>(null)
const ActiveSubmenu = ref<number | null>(null)

const GetActiveMenu = (menuId: number | null, subMenuId: number | null ) => {
  ActiveMainmenu.value = menuId
  ActiveSubmenu.value = subMenuId
  console.log(ActiveMainmenu.value, ActiveSubmenu.value)
}
menusData.value?.mainMenus.forEach(menu => {
  if (menu.menu_url === route.path) {
    GetActiveMenu(menu.menu_id, null)
  }
  menusData.value?.subMenus.forEach(subMenu => {
    if (subMenu.menu_url === route.path) {
      GetActiveMenu(subMenu.parent_id, subMenu.menu_id)
    }
  })
})
</script>
<template>
  <nav class="fixed top-0  left-0 h-screen flex flex-col border border-blue-200 bg-white z-50 transition-all duration-300" :class="LogoActive ? 'w-[50px]' : 'w-[250px]'">

    <div class="flex flex-row items-end justify-between mb-7 " :class="LogoActive ? 'px-1 py-0.5' : 'px-5 py-2 '">
      <img class="w-[45px] cursor-pointer" 
        @mouseover="LogoHover = true" @mouseleave="LogoHover = false" :src="LogoHover ? 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRFjMCSE996Z8SNXm4wjSLWxw1lFi8p28uocUgBwk3KuQ&s=10' : 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRDYI5Iquep92uL6WHGEgJGKS9UbC9tAoI8bsW1B_0ZwQ&s=10'" alt="logo"
        @click="LogoActive = !LogoActive"
      >
      <h1 class="text-[25px] font-bold transition-all duration-300" :class="LogoActive ? 'opacity-0' : 'opacity-1'">Akkarapon</h1>
    </div>

    <component :is="menu.menu_url == null ? 'div' : NuxtLink" :to="menu.menu_url" @click="menu.menu_url != null ? GetActiveMenu(menu.menu_id, null) : ''" v-for="menu in menusData?.mainMenus" :key="menu.menu_id" 
      class="relative group  items-center cursor-pointer" 
      @mouseover="hoveredMenuId = menu.menu_id" 
      @mouseleave="hoveredMenuId = null">
      <div class="flex flex-row group-hover:font-bold group-hover:-translate-y-0.5 group-hover:bg-gray-300 transition-all duration-300"
      :class="[ActiveMainmenu === menu.menu_id ? 'bg-gray-300 font-bold -translate-y-0.5 ' : '', LogoActive ? 'px-1 py-1 justify-center' : 'px-5 py-2 justify-between']">
        <Icon :name="`heroicons:${menu.menu_icon}-20-solid`" class="w-8 h-8 text-blue-400 " />
        <div class="flex flex-row justify-between w-[180px] items-end" :class="LogoActive ? 'hidden' : 'block'">
          <h1 class="text-[20px] transition-all duration-300 ml-2" :class="LogoActive ? 'opacity-0 hidden' : 'opacity-1 block'">{{ menu.name_menu }}</h1>
          <span v-if="menu.menu_url == null" class="pr-2">></span>
        </div>
      </div>
      <div v-if="menu.menu_url == null" class="absolute text-[20px] left-full flex flex-col top-0 right-10 bg-white text-black p-2 w-[250px] rounded-r-lg" :class="hoveredMenuId === menu.menu_id ? 'block' : 'hidden'">
        <NuxtLink :to="subMenu.menu_url!" @click="GetActiveMenu(subMenu.parent_id, subMenu.menu_id )" v-for="subMenu in menusData?.subMenus" :key="subMenu.menu_id" class="flex flex-row hover:bg-red-200 hover:text-black hover:translate-x-1 hover:-translate-y-1 transition-all duration-300" 
        :class="ActiveSubmenu === subMenu.menu_id ? 'bg-red-200 text-black translate-x-1 -translate-y-1 ' : ''">
          <Icon :name="`heroicons:${subMenu.menu_icon}-20-solid`" class="w-5 h-5 text-blue-400 mr-2" />
          <i>{{ subMenu.name_menu }}</i>
        </NuxtLink>
      </div>
    </component>
  </nav>
  <div class="transition-all duration-300" :class="LogoActive ? 'w-[50px]' : 'w-[250px]'"></div>
</template>