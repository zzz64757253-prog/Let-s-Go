<script setup lang="ts">
import { Confirm, Block, Loading, Notify } from 'notiflix/build/notiflix-aio';
const isMobile = useState<boolean>('isMobile')
const postPopup = ref<boolean>(false);
const delPopup = ref<boolean>(false);
const editPopup = ref<boolean>(false);
const pin = ref<string>('');
const post_uid = ref<string>('');
const title = ref<string>('');
const content = ref<string>('');
const username = ref<string>('');
const fileInput = ref<HTMLInputElement | null>(null);
const uploadedUrl = ref<string>('');
const uploadedUrl2 = ref<string>('');
const postsKey = ref<number>(0);
const namePost = ref<string>('');
const imageDelete = ref<boolean>(false);
const statusEdit = ref<boolean>(false);
type NotifyType = 'success' | 'failure' | 'warning' | 'info';
async function handleUpload(number: number) {
  if (number == 1){
    if (uploadedUrl.value && statusEdit.value == false){
      try {
        await $fetch('/server/delete', {
          method: 'DELETE',
          body: { url: uploadedUrl.value },
        })
        uploadedUrl.value = ''
        
      } catch (err: any) {
        alert('เกิดข้อผิดพลาด: ' + (err.data?.statusMessage || err.message))
        return
      }
    }
    const file = fileInput.value?.files?.[0] 
    if (!file) return alert('กรุณาเลือกไฟล์ก่อนครับ')

    const formData = new FormData()
    formData.append('file', file)

    Block.circle('#loadingImage', '', 
      {
        backgroundColor: 'transparent',
      }
    );
    try {
      const res = await $fetch<{url: string}>('/server/upload', {
        method: 'POST',
        body: formData,
      })
      if(statusEdit.value){
        uploadedUrl2.value = res.url
      }else{
        uploadedUrl.value = res.url
      }
    } catch (err: any) {
      alert('เกิดข้อผิดพลาด: ' + (err.data?.statusMessage || err.message))
    } finally {
      Block.remove('#loadingImage');
    }
  }else if (number == 2){

  }
}

async function handleSubmit() {
  try {
    // ใช้ $fetch ยิง POST ส่ง body ไปหา Server
    const result = await $fetch<{success:NotifyType, message:string}>('/server/content', {
      method: 'POST',
      body: {
        username: username.value,
        title: title.value,
        content: content.value,
        image_url: uploadedUrl.value,
        pin: pin.value,
      },
    })

    if (result.success == "failure"){
      Notify[result.success as NotifyType](result.message);
      return
    }
    username.value = ''
    title.value = ''
    content.value = ''
    uploadedUrl.value = ''
    pin.value = ''
    postPopup.value = false;
    Notify[result.success as NotifyType](result.message);
    postsKey.value++
  } catch (err: any) {
    if (uploadedUrl.value !== '') {
      await delImage(uploadedUrl.value)
    }
    Notify.failure(err.data?.statusMessage || err.message);
  }
}
  const delImage = async (url : string) => {
    await $fetch('/server/delete', {  
      method: 'DELETE',
      body: { url: url },
    })
  }
   const delImg2 = async (url : string) => {
    if (uploadedUrl.value === url) {
      uploadedUrl2.value = ""
      return
    }
    await $fetch('/server/delete', {
      method: 'DELETE',
      body: { url: url },
    })
  }
async function handleDelete() {
  delPopup.value = false
  Confirm.show(
    'ลบโพสต์',
    'คุณต้องการลบโพสต์ ใช่ไหม ?',
    'ใช่',
    'ไม่',
    async () => {
      postDelete()
    }, () => {
      delPopup.value = true
      return
    }

  )

  const postDelete = async () => {

    if (pin.value.length < 4 || post_uid.value === '') return Notify.failure('ไม่พบโพสต์ที่ต้องการลบ');
      Loading.circle('');
      try {
        const result = await $fetch<{success : NotifyType, message : string}>('/server/content', {
          method: 'DELETE',
          body: {
            post_uid: post_uid.value,
            pin: pin.value,
          },
        })
        post_uid.value = ''
        pin.value = ''
        delPopup.value = false
        Notify[result.success as NotifyType](result.message);
        if (uploadedUrl.value && result.success === "success"){
          await delImage(uploadedUrl.value)
          uploadedUrl.value = ''
        }
        Loading.remove();
        postsKey.value++
      } catch (err: any) {
        Loading.remove();
        Notify.failure(err.data?.statusMessage || err.message);
      }
  } 
}
const deleteImage = async () => {
  if (uploadedUrl.value === '') return
  console.log(uploadedUrl.value)
  const url = uploadedUrl.value
  
  Block.circle("#loadingImage", "")
  uploadedUrl.value = ""
  try {
    const result = await $fetch<{success: NotifyType; message: string }>('/server/delete', {
      method: 'DELETE',
      body: { url: url },
    })
    Block.remove('#loadingImage')
    Notify[result.success as NotifyType](result.message)
    uploadedUrl.value = ""
  }catch (err: any) {
    Block.remove('#loadingImage')
    Notify.failure(err.message)
  }
}
const editSubmit = async () => {
  editPopup.value = false
   Confirm.show(
    'แก้ไขโพสต์',
    'คุณต้องการแก้ไขโพสต์ ใช่ไหม ?',
    'ใช่',
    'ไม่',
    async () => {
      editPost()
    }, () => {
      editPopup.value = true
      return
    }
  )

  const editPost = async () => {
    const img = uploadedUrl.value
  
    uploadedUrl.value = uploadedUrl2.value
    try {
      const result = await $fetch<{success: NotifyType, message: string}>('/server/content', {
        method: 'PUT',
        body: {
          post_uid: post_uid.value,
          username: username.value,
          title: title.value,
          content: content.value,
          image_url: uploadedUrl.value,
          pin: pin.value,
        },
      })
      if (result.success === 'success' && img !== '' && uploadedUrl2.value !== img) {
        await delImage(img)
      }
      if (result.success === 'success'){
        username.value = ''
        title.value = ''
        content.value = ''
        uploadedUrl.value = ''
        pin.value = ''
        post_uid.value = ''
        editPopup.value = false;
        postsKey.value++
      }
      Notify[result.success as NotifyType](result.message);
    } catch (err: any) {
      Notify.failure(err.data?.statusMessage || err.message);
    }
  }
}

const hidePopup = async (number :number) => {
  if (number === 1 && uploadedUrl.value !== "") {
    await delImage(uploadedUrl.value)
  }else if (number === 2 && uploadedUrl2.value !== "") {
    await delImg2(uploadedUrl2.value)
  }
  pin.value = ''
  username.value = ''
  title.value = ''
  content.value = ''
  uploadedUrl.value = ''
  post_uid.value = ''
  statusEdit.value = false
  if (number === 1) postPopup.value = false
  if (number === 2) editPopup.value = false
  
}
const changeImg = async () => {
    await handleUpload(1)
}
function receiveFromChild(uid: string, name: string, imageUrl?: string) {
  post_uid.value = uid
  namePost.value = name
  uploadedUrl.value = imageUrl ?? ''
}
function receiveFromEdit(uid: string, name: string, postTitle: string, postContent: string, imageUrl?: string) {
  post_uid.value = uid
  username.value = name
  title.value = postTitle
  content.value = postContent
  uploadedUrl.value = imageUrl ?? ''
  if (imageUrl) uploadedUrl2.value = imageUrl
  statusEdit.value = true
}
</script>

<template>
  <div v-if="editPopup" class="fixed inset-0 flex justify-center items-center bg-black bg-opacity-50 z-[9999]">
    <form @submit.prevent="editSubmit" class="relative bg-white p-4 rounded-md w-[90%] max-w-[450px]"> 
      <div class="absolute top-1 right-1 flex justify-end">
        <Icon @click="hidePopup(2)" 
          name="heroicons:x-circle-20-solid" 
          class="w-6 h-6 text-gray-300 hover:text-red-500 cursor-pointer transition-all duration-300" />
        </div>
      <div class="flex flex-col gap-2">
        <div class="relative mt-4">
          <Icon name="heroicons:user-20-solid" class="absolute top-1 left-1 w-9 h-9 text-blue-400" />
          <input v-model="username" type="text" placeholder="ชื่อผู้โพสต์" class=" w-[90%] max-w-80 border border-gray-300 p-2 pl-12 rounded-md" />
        </div>
        <input v-model="title" type="text" placeholder="หัวข้อโพสต์" class="w-full border border-gray-300 p-2 rounded-md" />
        <div class="relative">
          <div @mouseover="imageDelete = true" @mouseleave="imageDelete = false" class="absolute bottom-5 right-5">
            <div v-if="uploadedUrl2" class="relative">
              <img :src="uploadedUrl2"  alt="Preview" class=" w-auto h-[65px] -rotate-12 transition-all duration-300" :class="imageDelete ? 'brightness-50 rotate-0' : ''" >
              <Icon v-if="imageDelete" @click="delImg2(uploadedUrl2)" name="heroicons:x-circle-20-solid"  class="absolute bottom-3 right-3 w-10 h-10 text-red-500 cursor-pointer" />
            </div>
          </div>
          <textarea v-model="content" placeholder="เนื้อหา" class="resize-none w-full h-32 border border-gray-300 p-2 rounded-md "></textarea>
          <div v-if="!uploadedUrl" class="absolute bottom-8 right-7 w-[60px] h-[60px]"><div id="loadingImage"></div></div>
        </div>
      </div>
      <div class="flex justify-between mt-2 gap-2">
        <label class="flex items-center gap-2 cursor-pointer">
          <Icon name="heroicons:photo-20-solid" class="text-green-400 h-9 w-9"/>
          <input type="file" ref="fileInput" @change="changeImg()" accept=".jpg, .jpeg, .png" class="hidden" />
        </label>
        <div class="relative">
          <Icon :class="pin.length == 4 ? 'opacity-1' : 'opacity-0'" name="heroicons:shield-check-20-solid" class="absolute top-1 right-1 w-9 h-9 text-green-400 transition-all duration-500" />
          <input v-model="pin"type="text" inputmode="numeric" maxlength="4" @input="pin = pin.replace(/\D/g, '')" class="border border-gray-300 p-2 rounded-md text-center w-[220px]" :class="pin.length > 0 ? 'tracking-[2em] indent-[2em] font-bold' : ''" placeholder="กรอก PIN ให้ถูกต้องเพื่อแก้ไข" />
        </div>
        <button type="submit" :disabled="pin.length < 4" class="bg-blue-400 w-[500px] rounded-md text-white disabled:bg-gray-300 disabled:text-gray-500">แก้ไขโพสต์</button>
      </div>
    </form>
  </div>

  <div v-if="delPopup" class="fixed inset-0 flex justify-center items-center bg-black bg-opacity-50 z-[9999]">
    <form  @submit.prevent="" class="relative bg-white p-4 rounded-md w-[90%] max-w-[400px]">
      <div class="absolute top-1 right-1 flex justify-end">
        <Icon @click="uploadedUrl = ''; pin = ''; delPopup = !delPopup " 
          name="heroicons:x-circle-20-solid" 
          class="w-6 h-6 text-gray-300 hover:text-red-500 cursor-pointer transition-all duration-300" />
        </div>
      <div class="bg-white p-4 rounded-md ">
        <h1 class="text-xl font-bold">ลบโพสต์ของ {{ namePost }}</h1>
        <p class="text-gray-600 text-[15px]">จำเป็นต้องใส่รหัส PIN เพิ่มลบโพสต์</p>
        <div class="flex flex-row justify-between items-center">
          <div class="relative ">
            <Icon :class="pin.length == 4 ? 'opacity-1' : 'opacity-0'" name="heroicons:shield-check-20-solid" class="absolute top-1 right-1 w-9 h-9 text-green-400 transition-all duration-500" />
            <input v-model="pin"type="text" inputmode="numeric" maxlength="4" @input="pin = pin.replace(/\D/g, '')" class="border border-gray-300 p-2 rounded-md text-center w-[220px] " :class="pin.length > 0 ? 'tracking-[2em] indent-[2em] font-bold' : ''" placeholder="PIN ใช้ตอนแก้ไข หรือลบโพสต์" />
          </div>
          <button type="submit" @click="handleDelete" :disabled="pin.length < 4" class="bg-red-400 px-4 py-2 rounded-md text-white disabled:bg-gray-300 disabled:text-gray-500">ลบ</button>
        </div>
      </div>
    </form>
  </div>


  <div v-if="postPopup" class="fixed inset-0 flex justify-center items-center bg-black bg-opacity-50 z-[9999]">
    <form @submit.prevent="handleSubmit" class="relative bg-white p-4 rounded-md w-[90%] max-w-[450px]"> 
      <div class="absolute top-1 right-1 flex justify-end">
        <Icon @click="hidePopup(1)" 
          name="heroicons:x-circle-20-solid" 
          class="w-6 h-6 text-gray-300 hover:text-red-500 cursor-pointer transition-all duration-300" />
        </div>
      <div class="flex flex-col gap-2">
        <div class="relative mt-4">
          <Icon name="heroicons:user-20-solid" class="absolute top-1 left-1 w-9 h-9 text-blue-400" />
          <input v-model="username" type="text" placeholder="ชื่อผู้โพสต์" class=" w-[90%] max-w-80 border border-gray-300 p-2 pl-12 rounded-md" />
        </div>
        <input v-model="title" type="text" placeholder="หัวข้อโพสต์" class="w-full border border-gray-300 p-2 rounded-md" />
        <div class="relative">
          <div @mouseover="imageDelete = true" @mouseleave="imageDelete = false" class="absolute bottom-5 right-5">
            <div class="relative">
              <img v-if="uploadedUrl" :src="uploadedUrl"  alt="Preview" class=" w-auto h-[65px] -rotate-12 transition-all duration-300" :class="imageDelete ? 'brightness-50 rotate-0' : ''" >
              <Icon v-if="imageDelete" @click="deleteImage()" name="heroicons:x-circle-20-solid"  class="absolute bottom-3 right-3 w-10 h-10 text-red-500 cursor-pointer" />
            </div>
          </div>
          <textarea v-model="content" placeholder="เนื้อหา" class="resize-none w-full h-32 border border-gray-300 p-2 rounded-md "></textarea>
          <div v-show="!uploadedUrl" class="absolute bottom-8 right-7 w-[60px] h-[60px]">
              <div id="loadingImage"></div>
          </div>
        </div>
      </div>
      <div class="flex justify-between mt-2 gap-2">
        <label class="flex items-center gap-2 cursor-pointer">
          <Icon name="heroicons:photo-20-solid" class="text-green-400 h-9 w-9"/>
          <input type="file" ref="fileInput" @change="handleUpload(1)" accept=".jpg, .jpeg, .png" class="hidden" />
        </label>
        <div class="relative">
          <Icon :class="pin.length == 4 ? 'opacity-1' : 'opacity-0'" name="heroicons:shield-check-20-solid" class="absolute top-1 right-1 w-9 h-9 text-green-400 transition-all duration-500" />
          <input v-model="pin"type="text" inputmode="numeric" maxlength="4" @input="pin = pin.replace(/\D/g, '')" class="border border-gray-300 p-2 rounded-md text-center w-[220px]" :class="pin.length > 0 ? 'tracking-[2em] indent-[2em] font-bold' : ''" placeholder="PIN ใช้ตอนแก้ไข หรือลบโพสต์" />
        </div>
        <button type="submit" :disabled="pin.length < 4" class="bg-blue-400 w-[500px] rounded-md text-white disabled:bg-gray-300 disabled:text-gray-500">โพสต์</button>
      </div>
    </form>
  </div>

  <div class="flex justify-center items-start w-full h-full">
  <div class="flex flex-col gap-4 w-full max-w-[800px] justify-center items-center rounded-md">
    <div :class="isMobile ? 'top-[60px]' : 'top-[0px]'" class="sticky z-50 flex w-full bg-white border border-gray-300 justify-center items-center p-4 flex-row gap-2 rounded-md">
      <p @click="postPopup = true" class="text-xs sm:text-sm md:text-lg text-gray-400 bg-gray-200 p-1 text-start pl-5 w-[500px] rounded-full cursor-pointer hover:bg-gray-300 ">โพสต์คอนเทนต์เกี่ยวกับคุณ...</p>
      <Icon name="heroicons:camera-20-solid" class="text-red-400 cursor-pointer h-9 w-9"/>
      <Icon name="heroicons:photo-20-solid" class="text-green-400 cursor-pointer h-9 w-9"/>
      <Icon name="heroicons:face-smile-20-solid" class="text-yellow-400 cursor-pointer h-9 w-9" />
    </div>
    <Posts :key="postsKey" 
    v-model:delPopup="delPopup" @send-uid="receiveFromChild" 
    v-model:editPopup="editPopup" @update="receiveFromEdit"
    />
  </div>
  </div>
</template>