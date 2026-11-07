<script setup lang="ts">
import { Block } from 'notiflix/build/notiflix-aio'

const openChat = ref<boolean>(false)
const message = ref('')
const messages = ref<{ content: string; role: string }[]>([])
const chatContainer = ref<HTMLDivElement | null>(null)
const isSending = ref(false)

interface Modeltype {
  modelname: string
}
interface ArrModel {
  models: Modeltype[]
}

const Needmodel = ref('gemini-3.5-flash-lite')
const { data: models, error } = await useFetch<ArrModel>('/server/model')

const scrollToBottom = async () => {
  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

watch(messages, () => { scrollToBottom() }, { deep: true })
watch(openChat, (isOpen) => { if (isOpen) scrollToBottom() })

const sendMessage = async () => {
  if (!message.value.trim() || isSending.value) return

  const userText = message.value.trim()
  message.value = ''
  messages.value.push({ content: userText, role: 'user' })
  isSending.value = true
  await nextTick()
  Block.dots('#thinking', 'Thinking')
  await scrollToBottom()

  try {
    const res = await $fetch<{ message: string }>('/server/chat', {
      method: 'POST',
      body: { messages: messages.value.slice(-5), model: Needmodel.value }
    })
    if (res.message) {
      messages.value.push({ content: res.message, role: 'assistant' })
    }
  } catch (err) {
    messages.value.push({
      content: 'เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง',
      role: 'assistant'
    })
    console.error('Error:', err)
  } finally {
    Block.remove('#thinking')
    isSending.value = false
    await scrollToBottom()
  }
}
</script>

<template>
  <!-- Chat popup: รองรับจอมือถือ / fixed ไม่ดันเนื้อหาหน้า -->
  <div
    v-if="openChat"
    role="dialog"
    aria-label="หน้าต่างแชต AI"
    class="fixed inset-x-3 bottom-[calc(12px+env(safe-area-inset-bottom))] z-[1000] flex h-[min(610px,calc(100dvh-80px))] flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl sm:inset-x-auto sm:bottom-24 sm:right-5 sm:h-[540px] sm:w-[380px]"
  >
    <!-- หัวแชตเรียบง่าย -->
    <div class="flex shrink-0 items-center justify-between gap-3 border-b border-gray-100 bg-white px-4 py-3">
      <div class="flex min-w-0 items-center gap-3">
        <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-50">
          <Icon name="heroicons:chat-bubble-left-right-20-solid" class="h-5 w-5 text-blue-400" />
        </span>
        <div class="min-w-0">
          <h2 class="text-sm font-bold text-gray-800">AI ผู้ช่วย</h2>
        </div>
      </div>
      <button
        type="button"
        aria-label="ปิดหน้าต่างแชต"
        class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-700"
        @click="openChat = false"
      >
        <Icon name="heroicons:x-mark-20-solid" class="h-5 w-5" />
      </button>
    </div>

    <!-- เลือกโมเดล แยกจากหัวเพื่อประหยัดพื้นที่และอ่านง่าย -->
    <div class="flex shrink-0 items-center gap-2 border-b border-gray-100 bg-gray-50/70 px-4 py-2">
      <Icon name="heroicons:cpu-chip-20-solid" class="h-4 w-4 shrink-0 text-blue-400" />
      <span class="shrink-0 text-xs text-gray-500">โมเดล</span>
      <select
        v-model="Needmodel"
        aria-label="เลือกโมเดล AI"
        class="min-w-0 flex-1 rounded-md border border-gray-200 bg-white px-2 py-1.5 text-xs text-gray-700 outline-none focus:border-blue-300 focus:ring-2 focus:ring-blue-100"
      >
        <option v-for="model in models?.models" :key="model.modelname" :value="model.modelname">
          {{ model.modelname }}
        </option>
      </select>
    </div>

    <!-- ข้อความ -->
    <div ref="chatContainer" class="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto bg-gray-50/60 p-4 scroll-smooth">
      <div v-if="messages.length === 0" class="my-auto flex flex-col items-center justify-center px-4 text-center">
        <span class="mb-3 flex h-14 w-14 items-center justify-center rounded-2xl border border-blue-100 bg-blue-50">
          <Icon name="heroicons:sparkles-20-solid" class="h-7 w-7 text-blue-400" />
        </span>
        <h3 class="text-sm font-semibold text-gray-700">เริ่มต้นบทสนทนา</h3>
        <p class="mt-1 text-xs leading-6 text-gray-400">พิมพ์ข้อความด้านล่างเพื่อพูดคุยกับ AI</p>
      </div>

      <div
        v-for="(msg, index) in messages"
        :key="index"
        class="max-w-[85%] break-words whitespace-pre-wrap rounded-2xl px-3.5 py-2.5 text-sm leading-6"
        :class="msg.role === 'user'
          ? 'self-end rounded-br-md bg-blue-400 text-white'
          : 'self-start rounded-bl-md border border-gray-200 bg-white text-gray-700'"
      >
        {{ msg.content }}
      </div>

      <!-- Notiflix loading เดิม -->
      <div v-show="isSending" id="thinking" class="relative min-h-[52px] w-36 self-start rounded-xl border border-gray-200 bg-white"></div>
    </div>

    <!-- ช่องพิมพ์ข้อความ -->
    <form class="flex shrink-0 items-end gap-2 border-t border-gray-100 bg-white p-3" @submit.prevent="sendMessage">
      <textarea
        v-model="message"
        rows="1"
        aria-label="ข้อความที่จะส่ง"
        placeholder="พิมพ์ข้อความ..."
        class="max-h-28 min-h-11 min-w-0 flex-1 resize-none rounded-xl border border-gray-200 bg-gray-50 px-3.5 py-3 text-sm leading-5 text-gray-700 outline-none transition-colors placeholder:text-gray-400 focus:border-blue-300 focus:bg-white focus:ring-2 focus:ring-blue-100"
        @keydown.enter.exact.prevent="sendMessage"
      ></textarea>
      <button
        type="submit"
        aria-label="ส่งข้อความ"
        :disabled="!message.trim() || isSending"
        class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-blue-400 text-white transition-colors hover:bg-blue-500 disabled:cursor-not-allowed disabled:opacity-40"
      >
        <Icon name="heroicons:paper-airplane-20-solid" class="h-5 w-5 -rotate-45" />
      </button>
    </form>
  </div>

  <!-- ปุ่มเปิดแชต -->
  <button
    v-else
    type="button"
    aria-label="เปิดแชต AI"
    class="fixed bottom-[calc(16px+env(safe-area-inset-bottom))] right-4 z-[1000] flex h-14 w-14 items-center justify-center rounded-2xl bg-blue-400 text-white shadow-lg shadow-blue-200 transition-all hover:-translate-y-1 hover:bg-blue-500 sm:right-5"
    @click="openChat = true"
  >
    <Icon name="heroicons:chat-bubble-oval-left-ellipsis-20-solid" class="h-7 w-7" />
  </button>
</template>
