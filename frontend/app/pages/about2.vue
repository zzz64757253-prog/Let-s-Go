<script setup lang="ts">

    const { data: users, pending, error } = await useFetch<any>('/api2/users2')
    console.log(users.value)
</script>
<template>
  <div v-for="user in users" :key="user.id">
    <i>{{ user.name }}</i>
  </div>
</template>