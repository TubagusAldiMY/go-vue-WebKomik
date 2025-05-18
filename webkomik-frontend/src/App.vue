<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { RouterLink, RouterView } from 'vue-router' // Import RouterView jika belum

const auth = useAuthStore()

onMounted(() => {
  auth.fetchUserSession() // Panggil saat komponen App di-mount
})

async function handleLogout() {
  await auth.logout()
}
</script>

<template>
  <div id="app-wrapper" class="min-h-screen bg-gray-50 text-gray-800">
    <header class="bg-white shadow-sm" v-if="auth.isAuthenticated">
      <nav class="container mx-auto px-4 py-3 flex justify-between items-center">
        <RouterLink :to="{ name: 'Home' }" class="text-xl font-semibold text-gray-800">WebKomik</RouterLink>
        <div>
          <span class="text-sm text-gray-600 mr-4">Welcome, {{ auth.userEmail }}</span>
          <button
              @click="handleLogout"
              class="text-sm font-medium text-indigo-600 hover:text-indigo-500 px-3 py-1 border border-indigo-600 rounded hover:bg-indigo-50 transition-colors"
          >
            Logout
          </button>
        </div>
      </nav>
    </header>

    <main class="container mx-auto p-4 mt-4">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
/* Style scoped khusus untuk App.vue jika diperlukan */
</style>