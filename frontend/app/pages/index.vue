<script setup>
// ดึงข้อมูลผ่าน Nuxt useFetch
const { data: users, pending, error } = await useFetch('/api/users')
</script>
<template>
  <div>
    HOME
    <div v-if="pending">Loading...</div>
    <div v-else-if="error">Error: {{ error.message }}</div>
    <ul v-else>
      <li v-for="user in users" :key="user.id">{{ user.name }}</li>
    </ul>
  </div>
  <ul>
    <li>Item 1</li>
    <li>Item 2</li>
    <li>Item 3</li>
  </ul>
</template>