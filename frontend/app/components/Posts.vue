<script setup lang="ts">
import { Block } from 'notiflix/build/notiflix-aio';
const thaiMonths = [
  '',
  'มกราคม',
  'กุมภาพันธ์',
  'มีนาคม',
  'เมษายน',
  'พฤษภาคม',
  'มิถุนายน',
  'กรกฎาคม',
  'สิงหาคม',
  'กันยายน',
  'ตุลาคม',
  'พฤศจิกายน',
  'ธันวาคม'
]
interface Post { 
  post_uid: string;
  username: string;
  title: string;
  content: string;
  image_url?: string;
  post_date: string
}

interface PostApiResponse {
  posts: Post[]
}
const { data: posts, pending, error } = await useFetch<PostApiResponse[]>('/server/content', {lazy: true})

onMounted(() => {
  watch(pending, (isPending) => {
    if (isPending) {
      Block.pulse('#loading', '')
    } else {
      Block.remove('#loading')
    }
  }, { immediate: true })
})

const delPopup = defineModel<boolean>('delPopup')
const emit = defineEmits<{
  'send-uid': [uid: string, name: string, imageUrl?: string]
  'update': [uid: string, name: string, title: string, content: string, imageUrl?: string]
}>()
function sendToParent(uid: string, name: string, imageUrl?: string) {
  emit('send-uid', uid, name, imageUrl) // ยิงค่ากลับไป
}
const editPopup = defineModel<boolean>('editPopup')
function sendToEdit(uid: string, name: string, title: string, content: string, imageUrl?: string) {
  emit('update', uid, name, title, content, imageUrl) // ยิงค่ากลับไป
}
const formatDate = (dateTime: string) => {
  const date = new
  Date().toLocaleDateString('en-EN')
  const datePart = dateTime.split("T")[0] ?? "";
  const [year, month, day] = datePart.split("-");
  const dateCheck = `${month}/${day}/${year}`
  if (date === dateCheck) {
    return 'วันนี้'
  }else{
    return `${day} ${thaiMonths[Number(month)]}`
  }
  
}
</script>

<template>
    <div id="loading" class="w-[100%]">
        <div v-if="posts?.posts" v-for="post in posts.posts" :key="post.post_uid" class="flex flex-col w-[100%] bg-white border border-gray-300 justify-start items-start p-4 flex-row gap-2 rounded-md mb-2">
            <div class="flex flex-row justify-between w-full">
                <p class="text-gray-600 text-[18px] items-center flex"><Icon name="heroicons:user-20-solid" class="text-blue-400 cursor-pointer h-6 w-6" /> : {{ post.username }}</p>
              <div class="flex flex-row gap-5">
                  <Icon @click="sendToEdit(post.post_uid, post.username, post.title, post.content, post.image_url); editPopup = !editPopup" name="heroicons:pencil-square-20-solid" class="text-blue-400 cursor-pointer h-6 w-6" />
                  <Icon @click="sendToParent(post.post_uid, post.username, post.image_url ); delPopup = !delPopup" name="heroicons:trash-20-solid" class="text-red-400 cursor-pointer h-6 w-6" />
              </div>
            </div>
            <h1 class="text-xl font-bold">{{ post.title }}</h1>
            <p class="text-gray-600 text-[15px]">{{ post.content }}</p>
            <img v-if="post.image_url" :src="post.image_url" :alt="post.username" class="w-full">
            <div class="flex flex-row justify-end w-[100%]">
              <p class="text-gray-600 text-[13px] ml-2">{{ formatDate(post.post_date) }}</p>
            </div>
        </div>
    </div>
</template>