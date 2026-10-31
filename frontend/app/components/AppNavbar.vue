<script setup lang="ts">
import { NuxtLink } from '#components'
const route = useRoute()
console.log(route.path) 
interface MenuType {
  menu_id: number
  name_menu: string
  menu_url: string
  menu_icon: string
}
const { data: menus } = await useFetch<MenuType[]>('/api/mobilemenu')
console.log(menus.value)
</script>
<template>
  <nav class="fixed top-0 w-screen h-[60px] border border-black bg-white z-50">
    <div class="flex justify-between items-center border border-blue-200 h-[100%] px-5 py-0">
        <div v-for="menu in menus" :key="menu.menu_id" class="h-full w-[50px]" >
            <NuxtLink :to="menu.menu_url" class="relative h-full flex justify-center items-center">
                <Icon :name="`heroicons:${menu.menu_icon}-20-solid`" 
                    class="w-8 h-8 text-blue-400" 
                    :class="menu.menu_url === route.path ? 'text-blue-400' : 'text-gray-400'"  />
                <span :class="menu.menu_url === route.path ? 'w-full border-blue-400' : 'w-[0px] border-white'" class="absolute bottom-0 left-0 border-2 border-blue-400  transition-all duration-300"></span>
            </NuxtLink>
        </div>
    </div>
  </nav>
  <div class="h-[60px]"></div>
</template>