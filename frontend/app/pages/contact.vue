<script setup>
// 1. ตัวแปรเก็บข้อความที่พิมพ์ค้นหา
const searchQuery = ref('')

// 2. ข้อมูลตั้งต้น (ตัวอย่างโพสต์)
const posts = ref([
  { id: 1, title: 'KhunZee', author: 'KhunZee' },
  { id: 2, title: 'วิธีจัดการ CSS ใน Firefox', author: 'DevA' },
  { id: 3, title: 'ทำระบบค้นหาด้วย Vue 3', author: 'KhunZee' },
  { id: 4, title: 'Tailwind CSS Tips & Tricks', author: 'Admin' }
])

// 3. Computed สำหรับกรองข้อมูลตามคำค้นหา
const filteredPosts = computed(() => {
  // ตัดช่องว่างหน้า-หลัง
  const keyword = searchQuery.value.trim().toLowerCase()

  // ถ้ายังไม่ได้พิมพ์อะไร ให้แสดงข้อมูลทั้งหมด
  if (!keyword) return posts.value

  return posts.value.filter(post => {
    // แปลงข้อความที่จะค้นหาให้เป็นพิมพ์เล็กทั้งหมดเพื่อไม่ให้ติด Case-sensitive (A != a)
    const matchTitle = post.title.toLowerCase().includes(keyword)
    const matchAuthor = post.author.toLowerCase().includes(keyword) // ค้นจากชื่อผู้แต่งด้วย (ถ้าต้องการ)
    
    return matchTitle || matchAuthor
  })
})
</script>

<template>
  <div class="max-w-md mx-auto p-4">
    <!-- ช่อง Input ค้นหา -->
    <div class="relative mb-4">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="ค้นหาโพสต์..."
        class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <!-- ปุ่มเคลียร์คำค้นหาเมื่อมีข้อความ -->
      <button
        v-if="searchQuery"
        @click="searchQuery = ''"
        class="absolute right-3 top-2.5 text-gray-400 hover:text-gray-600"
      >
        ✕
      </button>
    </div>

    <!-- ลูปแสดงผลรายการที่ค้นพบ -->
    <div class="space-y-2">
      <div
        v-for="post in filteredPosts"
        :key="post.id"
        class="p-3 bg-white border rounded-lg shadow-sm hover:shadow transition"
      >
        <h3 class="font-semibold text-gray-800">{{ post.title }}</h3>
        <p class="text-sm text-gray-500">ผู้เขียน: {{ post.author }}</p>
      </div>

      <!-- กรณีค้นหาแล้วไม่เจออะไรเลย -->
      <div
        v-if="filteredPosts.length === 0"
        class="text-center py-8 text-gray-400"
      >
        ไม่พบข้อมูลที่ตรงกับ "{{ searchQuery }}"
      </div>
    </div>
  </div>
</template>