<script setup lang="ts">
import { Block } from 'notiflix/build/notiflix-aio';

const openChat = ref<boolean>(false);
const message = ref('');
const messages = ref<{ content: string; role: string }[]>([]);

const chatContainer = ref<HTMLDivElement | null>(null);
interface Modeltype {
  modelname: string
}
interface ArrModel {
  models: Modeltype[]
}
const Needmodel = ref("gemini-3.5-flash-lite");

const { data: models, error } = await useFetch<ArrModel>('/server/model')
const scrollToBottom = async () => {
  await nextTick(); // รอให้ Vue วาดข้อความใหม่ลง DOM ก่อน
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};
watch(messages, () => {
  scrollToBottom();
},{ deep: true });
watch(openChat, (isOpen) => {
  if (isOpen) scrollToBottom();
});

const sendMessage = async () => {
  if (!message.value) return
  Block.dots('#thinking', 'Thinking')
  try {
    const obj = {
      content: message.value,
      role: "user"
    };
    messages.value.push(obj);
    message.value = '';
    console.log(Needmodel.value)
    const res = await $fetch<{ message: string }>('/server/chat', {
      method: 'POST',
      body: { messages: messages.value.slice(-5), model: Needmodel.value }
    })
    if (res.message){
      const obj2 = {
        content: res.message,
        role: "assistant"
      };
      messages.value.push(obj2);
      Block.remove('#thinking')
    }
  } catch (err) {
    const obj3 = {
      content: 'เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง',
      role: "assistant"
    };
    messages.value.push(obj3);
    Block.remove('#thinking')
    console.error('Error:', err)
  }
}
</script>

<template>
  <div :class="openChat == true ? 'block' : 'hidden'" class="border border-gray-300 rounded-md w-80 h-[420px] flex flex-col overflow-hidden">
    <div class="bg-blue-400 p-2 flex flex-row justify-between items-center">
      <div class="flex"> <h1 class="bg-blue-400 text-white text-xl font-bold">คุยกับ</h1>
        <select v-model="Needmodel" class="bg-blue-400 text-white text-sm font-bold">
          <option class="bg-gray-400 text-white text-xl font-bold" 
            v-for="model in models?.models"
            :key="model.modelname"
            :value="model.modelname"
          >
            {{ model.modelname }}
          </option>
        </select>
      </div>
      <Icon @click="openChat = !openChat" name="heroicons:x-circle-20-solid" class="w-6 h-6 text-white hover:text-red-500 cursor-pointer transition-all duration-300" />
    </div>
    <div 
      ref="chatContainer" 
      class="flex-1 p-2 flex flex-col overflow-y-auto scroll-smooth bg-white"
    >
      <p v-for="(msg, index) in messages" 
        :key="index"
        class=" bg-blue-400 w-fit max-w-[70%] break-words p-2 my-0.5" 
        :class="msg.role === 'user' ? 'self-end rounded-bl-lg rounded-t-lg' : 'self-start rounded-br-lg rounded-t-lg'"
      >
        {{ msg.content }}
    </p>
    <div id="thinking"></div>
    <div class="h-[50px] shrink-0"></div>
    </div>
    <div class="relative  h-10">
      <div @click="sendMessage" class="absolute right-2 bottom-1 w-7 h-7 bg-blue-500 rounded-full flex items-center justify-center hover:-rotate-45 transition-all duration-300">
        <Icon name="heroicons:paper-airplane-20-solid" class="w-6 h-6 text-white -rotate-45 cursor-pointer " />
      </div>
      <textarea type="text" 
        v-model="message" 
        @keydown.enter.exact.prevent="sendMessage" 
        @keydown.shift.enter.exact=""
        class="w-full border border-gray-300 h-10 rounded pr-10 pl-2 text-xl pt-1.5 bg-white" ></textarea>
    </div>
  </div>
  <div @click="openChat = !openChat" :class="openChat == true ? 'hidden' : 'block'"
    class="fixed bottom-4 right-4 w-14 h-14 bg-blue-500 rounded-full flex items-center justify-center cursor-pointer z-50 hover:bg-blue-600 hover:-rotate-45 transition-all duration-300">
    <Icon name="heroicons:chat-bubble-oval-left-ellipsis-20-solid" class="w-8 h-8 text-white " />
  </div>
</template>
