<script setup lang="ts">
import { Block } from 'notiflix/build/notiflix-aio'

const thaiMonths = [
  '',
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน',
  'พฤษภาคม', 'มิถุนายน', 'กรกฎาคม', 'สิงหาคม',
  'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
]

interface Post {
  post_uid: string
  username: string
  title: string
  content: string
  image_url?: string
  post_date: string
}

interface PostApiResponse {
  posts: Post[]
}

// API ส่งกลับเป็น { posts: [...] } จึงใช้ PostApiResponse (ไม่ใช่ PostApiResponse[])
const { data: posts, pending, error } = await useFetch<PostApiResponse>(
  '/server/content',
  { lazy: true }
)

onMounted(() => {
  watch(pending, (isPending) => {
    if (isPending) {
      Block.pulse('#loading', '', { backgroundColor: 'transparent' })
    } else {
      Block.remove('#loading')
    }
  }, { immediate: true })
})

onBeforeUnmount(() => Block.remove('#loading'))

const delPopup = defineModel<boolean>('delPopup')
const editPopup = defineModel<boolean>('editPopup')

const emit = defineEmits<{
  'send-uid': [uid: string, name: string, imageUrl?: string]
  'update': [uid: string, name: string, title: string, content: string, imageUrl?: string]
}>()

function sendToParent(uid: string, name: string, imageUrl?: string) {
  emit('send-uid', uid, name, imageUrl)
}

function sendToEdit(uid: string, name: string, title: string, content: string, imageUrl?: string) {
  emit('update', uid, name, title, content, imageUrl)
}

const formatDate = (dateTime: string) => {
  const date = new Date().toLocaleDateString('en-EN')
  const datePart = dateTime.split('T')[0] ?? ''
  const [year, month, day] = datePart.split('-')
  const dateCheck = `${month}/${day}/${year}`
  if (date === dateCheck) {
    return 'วันนี้'
  }
  return `${day} ${thaiMonths[Number(month)]}`
}
</script>

<template>
  <section id="loading" class="w-full min-w-0 text-gray-700">
    <!-- ส่วนหัวกะทัดรัด ไม่กินพื้นที่บนมือถือ -->
    <div class="mb-4 flex items-center justify-between gap-3 rounded-xl border border-gray-200 bg-white px-4 py-3 sm:px-5">
      <div class="flex min-w-0 items-center gap-3">
        <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-blue-50">
          <Icon name="heroicons:document-text-20-solid" class="h-5 w-5 text-blue-400" />
        </span>
        <div class="min-w-0">
          <h2 class="text-base font-bold text-gray-800 sm:text-lg">โพสต์ล่าสุด</h2>
          <p class="text-xs text-gray-400">ข่าวสารและเนื้อหาที่แบ่งปัน</p>
        </div>
      </div>
      <span v-if="posts?.posts" class="shrink-0 rounded-full bg-gray-100 px-3 py-1 text-xs text-gray-500">
        {{ posts.posts.length }} โพสต์
      </span>
    </div>

    <div v-if="pending && !posts?.posts" class="space-y-3" aria-label="กำลังโหลดโพสต์">
      <div v-for="n in 2" :key="n" class="animate-pulse rounded-xl border border-gray-200 bg-white p-4 sm:p-5">
        <div class="mb-5 flex items-center gap-3">
          <div class="h-10 w-10 rounded-full bg-gray-100"></div>
          <div class="space-y-2"><div class="h-3 w-24 rounded bg-gray-100"></div><div class="h-2 w-14 rounded bg-gray-100"></div></div>
        </div>
        <div class="mb-3 h-4 w-2/3 rounded bg-gray-100"></div>
        <div class="mb-2 h-3 rounded bg-gray-100"></div>
        <div class="h-3 w-4/5 rounded bg-gray-100"></div>
      </div>
    </div>

    <div v-else-if="error" class="rounded-xl border border-red-100 bg-white px-4 py-8 text-center">
      <Icon name="heroicons:exclamation-circle-20-solid" class="mx-auto mb-2 h-7 w-7 text-red-400" />
      <p class="text-sm text-gray-600">ไม่สามารถโหลดโพสต์ได้ กรุณาลองใหม่อีกครั้ง</p>
    </div>

    <div v-else-if="posts?.posts?.length" class="space-y-3">
      <article
        v-for="post in posts.posts"
        :key="post.post_uid"
        class="min-w-0 overflow-hidden rounded-xl border border-gray-200 bg-white p-4 transition-colors hover:border-blue-200 sm:p-5"
      >
        <!-- ผู้เขียน + ปุ่มจัดการ -->
        <div class="mb-4 flex items-start justify-between gap-2">
          <div class="flex min-w-0 items-center gap-3">
            <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-blue-50">
              <Icon name="heroicons:user-20-solid" class="h-5 w-5 text-blue-400" />
            </div>
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold text-gray-800">{{ post.username }}</p>
              <p class="mt-0.5 text-xs text-gray-400">{{ formatDate(post.post_date) }}</p>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-1">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-blue-50 hover:text-blue-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-300"
              :aria-label="`แก้ไขโพสต์ ${post.title}`"
              title="แก้ไขโพสต์"
              @click="sendToEdit(post.post_uid, post.username, post.title, post.content, post.image_url); editPopup = !editPopup"
            >
              <Icon name="heroicons:pencil-square-20-solid" class="h-5 w-5" />
            </button>
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-red-50 hover:text-red-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-300"
              :aria-label="`ลบโพสต์ ${post.title}`"
              title="ลบโพสต์"
              @click="sendToParent(post.post_uid, post.username, post.image_url); delPopup = !delPopup"
            >
              <Icon name="heroicons:trash-20-solid" class="h-5 w-5" />
            </button>
          </div>
        </div>

        <h3 class="mb-2 break-words text-base font-bold leading-7 text-gray-800 sm:text-lg">{{ post.title }}</h3>
        <p class="break-words whitespace-pre-wrap text-sm leading-7 text-gray-600">{{ post.content }}</p>

        <div v-if="post.image_url" class="mt-4 overflow-hidden rounded-lg border border-gray-100 bg-gray-50">
          <img
            :src="post.image_url"
            :alt="post.title || post.username"
            class="block max-h-[450px] w-full object-contain"
            loading="lazy"
          />
        </div>
      </article>
    </div>

    <div v-else-if="!pending" class="rounded-xl border border-gray-200 bg-white px-4 py-10 text-center">
      <Icon name="heroicons:document-text" class="mx-auto mb-2 h-8 w-8 text-gray-300" />
      <p class="text-sm text-gray-500">ยังไม่มีโพสต์ในขณะนี้</p>
    </div>
  </section>
</template>
