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
  <div class="min-h-dvh w-full bg-gray-50 px-3 py-3 text-gray-700 sm:px-5 sm:py-6">
    <main class="mx-auto flex w-full max-w-[800px] flex-col gap-3 sm:gap-4">

      <!-- แถบสร้างโพสต์ : กะทัดรัดบนมือถือ -->
      <section
        :class="isMobile ? 'top-[60px]' : 'top-0'"
        class="sticky z-40 rounded-lg border border-gray-200 bg-white p-3 shadow-sm sm:p-4"
        aria-label="สร้างโพสต์ใหม่"
      >
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-blue-50">
            <Icon name="heroicons:user-20-solid" class="h-5 w-5 text-blue-400" />
          </div>

          <button
            type="button"
            @click="postPopup = true"
            class="min-w-0 flex-1 rounded-full border border-gray-200 bg-gray-50 px-4 py-2.5 text-left text-sm text-gray-400 transition-colors hover:border-blue-200 hover:bg-blue-50/50"
          >
            โพสต์คอนเทนต์เกี่ยวกับคุณ...
          </button>

          <button
            type="button"
            @click="postPopup = true"
            class="flex h-10 shrink-0 items-center justify-center gap-1.5 rounded-md bg-blue-400 px-3 text-sm font-medium text-white transition-colors hover:bg-blue-500 sm:px-4"
            aria-label="สร้างโพสต์"
          >
            <Icon name="heroicons:plus-20-solid" class="h-5 w-5" />
            <span class="hidden sm:inline">สร้างโพสต์</span>
          </button>
        </div>
      </section>

      <!-- ใช้ Component และ Events เดิม -->
      <Posts
        :key="postsKey"
        v-model:delPopup="delPopup"
        @send-uid="receiveFromChild"
        v-model:editPopup="editPopup"
        @update="receiveFromEdit"
      />
    </main>
  </div>

  <!-- POPUP : สร้างโพสต์ -->
  <div
    v-if="postPopup"
    class="fixed inset-0 z-[9999] flex items-center justify-center bg-gray-950/50 p-3 backdrop-blur-[2px] sm:p-5"
  >
    <form
      @submit.prevent="handleSubmit"
      class="flex w-full max-w-[500px] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl"
      style="max-height: calc(100dvh - 24px)"
      role="dialog"
      aria-modal="true"
      aria-labelledby="create-post-title"
    >
      <div class="flex shrink-0 items-center justify-between border-b border-gray-100 px-4 py-3.5 sm:px-5">
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-md bg-blue-50">
            <Icon name="heroicons:pencil-square-20-solid" class="h-5 w-5 text-blue-400" />
          </div>
          <div>
            <h2 id="create-post-title" class="text-base font-bold text-gray-800">สร้างโพสต์</h2>
            <p class="text-xs text-gray-400">แบ่งปันสิ่งที่คุณอยากเล่า</p>
          </div>
        </div>
        <button
          type="button"
          @click="hidePopup(1)"
          class="flex h-9 w-9 items-center justify-center rounded-md text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-700"
          aria-label="ปิดหน้าต่าง"
        >
          <Icon name="heroicons:x-mark-20-solid" class="h-5 w-5" />
        </button>
      </div>

      <div class="min-h-0 flex-1 space-y-4 overflow-y-auto px-4 py-4 sm:px-5">
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-gray-600">ชื่อผู้โพสต์</label>
          <input
            v-model="username"
            type="text"
            placeholder="ชื่อของคุณ"
            class="w-full rounded-md border border-gray-200 bg-white px-3 py-2.5 text-sm outline-none transition-colors placeholder:text-gray-400 focus:border-blue-400 focus:ring-2 focus:ring-blue-100"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-semibold text-gray-600">หัวข้อโพสต์</label>
          <input
            v-model="title"
            type="text"
            placeholder="หัวข้อที่ต้องการโพสต์"
            class="w-full rounded-md border border-gray-200 bg-white px-3 py-2.5 text-sm outline-none transition-colors placeholder:text-gray-400 focus:border-blue-400 focus:ring-2 focus:ring-blue-100"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-semibold text-gray-600">เนื้อหา</label>
          <textarea
            v-model="content"
            rows="5"
            placeholder="เขียนเนื้อหาของคุณที่นี่..."
            class="w-full resize-y rounded-md border border-gray-200 bg-white px-3 py-2.5 text-sm leading-6 outline-none transition-colors placeholder:text-gray-400 focus:border-blue-400 focus:ring-2 focus:ring-blue-100"
          ></textarea>
        </div>

        <!-- รูปภาพและสถานะอัปโหลด -->
        <div class="rounded-lg border border-dashed border-gray-200 bg-gray-50/60 p-3">
          <div class="flex items-center justify-between gap-3">
            <div class="flex min-w-0 items-center gap-2">
              <Icon name="heroicons:photo-20-solid" class="h-5 w-5 text-blue-400" />
              <span class="text-sm font-medium text-gray-600">รูปภาพประกอบ</span>
            </div>
            <label class="inline-flex shrink-0 cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-semibold text-gray-600 transition-colors hover:border-blue-300 hover:text-blue-500">
              <Icon name="heroicons:arrow-up-tray-20-solid" class="h-4 w-4" />
              เลือกรูปภาพ
              <input ref="fileInput" type="file" accept=".jpg,.jpeg,.png" class="hidden" @change="handleUpload(1)" />
            </label>
          </div>

          <div class="relative mt-2 min-h-[36px]">
            <div id="loadingImage" class="absolute right-0 top-0 h-9 w-9"></div>
            <div v-if="uploadedUrl" class="relative mt-2 overflow-hidden rounded-md border border-gray-200 bg-white">
              <img :src="uploadedUrl" alt="รูปตัวอย่างก่อนโพสต์" class="max-h-48 w-full object-contain" />
              <button
                type="button"
                @click="deleteImage()"
                class="absolute right-2 top-2 flex h-8 w-8 items-center justify-center rounded-md bg-white/95 text-red-500 shadow-sm transition-colors hover:bg-red-50"
                aria-label="ลบรูปภาพ"
              >
                <Icon name="heroicons:x-mark-20-solid" class="h-5 w-5" />
              </button>
            </div>
            <p v-else class="pr-10 text-xs text-gray-400">รองรับไฟล์ JPG และ PNG</p>
          </div>
        </div>
      </div>

      <!-- PIN และปุ่มจัดการ : ไม่เบียดกันบนจอเล็ก -->
      <div class="shrink-0 space-y-3 border-t border-gray-100 bg-white px-4 py-3.5 sm:px-5">
        <div class="flex items-center gap-3">
          <label for="create-post-pin" class="shrink-0 text-xs font-semibold text-gray-600">PIN 4 หลัก</label>
          <div class="relative min-w-0 flex-1">
            <input
              id="create-post-pin"
              v-model="pin"
              type="password"
              inputmode="numeric"
              maxlength="4"
              @input="pin = pin.replace(/\D/g, '')"
              placeholder="กรอกรหัส PIN"
              class="w-full rounded-md border border-gray-200 bg-white px-3 py-2.5 pr-10 text-center text-sm tracking-[0.3em] outline-none transition-colors focus:border-blue-400 focus:ring-2 focus:ring-blue-100"
            />
            <Icon
              v-if="pin.length === 4"
              name="heroicons:shield-check-20-solid"
              class="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-green-500"
            />
          </div>
        </div>
        <div class="flex gap-2">
          <button type="button" @click="hidePopup(1)" class="flex-1 rounded-md border border-gray-200 bg-white px-4 py-2.5 text-sm font-semibold text-gray-600 hover:bg-gray-50">ยกเลิก</button>
          <button type="submit" :disabled="pin.length < 4" class="flex-1 rounded-md bg-blue-400 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-blue-500 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400">เผยแพร่โพสต์</button>
        </div>
      </div>
    </form>
  </div>

  <!-- POPUP : แก้ไขโพสต์ -->
  <div
    v-if="editPopup"
    class="fixed inset-0 z-[9999] flex items-center justify-center bg-gray-950/50 p-3 backdrop-blur-[2px] sm:p-5"
  >
    <form
      @submit.prevent="editSubmit"
      class="flex w-full max-w-[500px] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl"
      style="max-height: calc(100dvh - 24px)"
      role="dialog"
      aria-modal="true"
      aria-labelledby="edit-post-title"
    >
      <div class="flex shrink-0 items-center justify-between border-b border-gray-100 px-4 py-3.5 sm:px-5">
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-md bg-blue-50">
            <Icon name="heroicons:pencil-square-20-solid" class="h-5 w-5 text-blue-400" />
          </div>
          <div>
            <h2 id="edit-post-title" class="text-base font-bold text-gray-800">แก้ไขโพสต์</h2>
            <p class="text-xs text-gray-400">อัปเดตข้อความและรูปภาพ</p>
          </div>
        </div>
        <button type="button" @click="hidePopup(2)" class="flex h-9 w-9 items-center justify-center rounded-md text-gray-400 hover:bg-gray-100 hover:text-gray-700" aria-label="ปิดหน้าต่าง">
          <Icon name="heroicons:x-mark-20-solid" class="h-5 w-5" />
        </button>
      </div>

      <div class="min-h-0 flex-1 space-y-4 overflow-y-auto px-4 py-4 sm:px-5">
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-gray-600">ชื่อผู้โพสต์</label>
          <input v-model="username" type="text" placeholder="ชื่อของคุณ" class="w-full rounded-md border border-gray-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100" />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-gray-600">หัวข้อโพสต์</label>
          <input v-model="title" type="text" placeholder="หัวข้อโพสต์" class="w-full rounded-md border border-gray-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100" />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-gray-600">เนื้อหา</label>
          <textarea v-model="content" rows="5" placeholder="แก้ไขเนื้อหาของโพสต์..." class="w-full resize-y rounded-md border border-gray-200 bg-white px-3 py-2.5 text-sm leading-6 outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100"></textarea>
        </div>

        <div class="rounded-lg border border-dashed border-gray-200 bg-gray-50/60 p-3">
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center gap-2 text-sm font-medium text-gray-600">
              <Icon name="heroicons:photo-20-solid" class="h-5 w-5 text-blue-400" />
              รูปภาพประกอบ
            </div>
            <label class="inline-flex shrink-0 cursor-pointer items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-2 text-xs font-semibold text-gray-600 hover:border-blue-300 hover:text-blue-500">
              <Icon name="heroicons:arrow-up-tray-20-solid" class="h-4 w-4" />
              เปลี่ยนรูป
              <input ref="fileInput" type="file" accept=".jpg,.jpeg,.png" class="hidden" @change="changeImg()" />
            </label>
          </div>
          <div class="relative mt-2 min-h-[36px]">
            <div id="loadingImage" class="absolute right-0 top-0 h-9 w-9"></div>
            <div v-if="uploadedUrl2" class="relative mt-2 overflow-hidden rounded-md border border-gray-200 bg-white">
              <img :src="uploadedUrl2" alt="รูปภาพประกอบโพสต์" class="max-h-48 w-full object-contain" />
              <button
                type="button"
                @click="delImg2(uploadedUrl2)"
                class="absolute right-2 top-2 flex h-8 w-8 items-center justify-center rounded-md bg-white/95 text-red-500 shadow-sm hover:bg-red-50"
                aria-label="ลบรูปภาพที่เลือก"
              >
                <Icon name="heroicons:x-mark-20-solid" class="h-5 w-5" />
              </button>
            </div>
            <p v-else class="pr-10 text-xs text-gray-400">ยังไม่ได้เลือกรูปภาพ</p>
          </div>
        </div>
      </div>

      <div class="shrink-0 space-y-3 border-t border-gray-100 bg-white px-4 py-3.5 sm:px-5">
        <div class="flex items-center gap-3">
          <label for="edit-post-pin" class="shrink-0 text-xs font-semibold text-gray-600">PIN 4 หลัก</label>
          <div class="relative min-w-0 flex-1">
            <input
              id="edit-post-pin"
              v-model="pin"
              type="password"
              inputmode="numeric"
              maxlength="4"
              @input="pin = pin.replace(/\D/g, '')"
              placeholder="PIN สำหรับยืนยันการแก้ไข"
              class="w-full rounded-md border border-gray-200 bg-white px-3 py-2.5 pr-10 text-center text-sm tracking-[0.3em] outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100"
            />
            <Icon v-if="pin.length === 4" name="heroicons:shield-check-20-solid" class="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-green-500" />
          </div>
        </div>
        <div class="flex gap-2">
          <button type="button" @click="hidePopup(2)" class="flex-1 rounded-md border border-gray-200 bg-white px-4 py-2.5 text-sm font-semibold text-gray-600 hover:bg-gray-50">ยกเลิก</button>
          <button type="submit" :disabled="pin.length < 4" class="flex-1 rounded-md bg-blue-400 px-4 py-2.5 text-sm font-semibold text-white hover:bg-blue-500 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400">บันทึกการแก้ไข</button>
        </div>
      </div>
    </form>
  </div>

  <!-- POPUP : ลบโพสต์ -->
  <div
    v-if="delPopup"
    class="fixed inset-0 z-[9999] flex items-center justify-center bg-gray-950/50 p-3 backdrop-blur-[2px]"
  >
    <form
      @submit.prevent="handleDelete"
      class="w-full max-w-[420px] overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl"
      role="dialog"
      aria-modal="true"
      aria-labelledby="delete-post-title"
    >
      <div class="flex items-center justify-between border-b border-gray-100 px-4 py-3.5 sm:px-5">
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-md bg-red-50">
            <Icon name="heroicons:trash-20-solid" class="h-5 w-5 text-red-500" />
          </div>
          <h2 id="delete-post-title" class="text-base font-bold text-gray-800">ลบโพสต์</h2>
        </div>
        <button type="button" @click="uploadedUrl = ''; pin = ''; delPopup = false" class="flex h-9 w-9 items-center justify-center rounded-md text-gray-400 hover:bg-gray-100 hover:text-gray-700" aria-label="ปิดหน้าต่าง">
          <Icon name="heroicons:x-mark-20-solid" class="h-5 w-5" />
        </button>
      </div>

      <div class="space-y-4 px-4 py-5 sm:px-5">
        <div class="rounded-md border border-red-100 bg-red-50/60 p-3">
          <p class="text-sm font-semibold text-gray-800 break-words">โพสต์ของ {{ namePost }}</p>
          <p class="mt-1 text-xs leading-5 text-gray-500">กรอก PIN 4 หลักเพื่อยืนยันการลบโพสต์</p>
        </div>
        <div class="relative">
          <input
            v-model="pin"
            type="password"
            inputmode="numeric"
            maxlength="4"
            @input="pin = pin.replace(/\D/g, '')"
            placeholder="PIN 4 หลัก"
            class="w-full rounded-md border border-gray-200 bg-white px-3 py-3 pr-10 text-center text-sm tracking-[0.3em] outline-none focus:border-red-400 focus:ring-2 focus:ring-red-100"
            aria-label="PIN สำหรับลบโพสต์"
          />
          <Icon v-if="pin.length === 4" name="heroicons:shield-check-20-solid" class="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-green-500" />
        </div>
        <div class="flex gap-2">
          <button type="button" @click="uploadedUrl = ''; pin = ''; delPopup = false" class="flex-1 rounded-md border border-gray-200 px-4 py-2.5 text-sm font-semibold text-gray-600 hover:bg-gray-50">ยกเลิก</button>
          <button type="submit" :disabled="pin.length < 4" class="flex-1 rounded-md bg-red-500 px-4 py-2.5 text-sm font-semibold text-white hover:bg-red-600 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400">ยืนยันลบโพสต์</button>
        </div>
      </div>
    </form>
  </div>
</template>
