<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/authStore' // Pastikan path ini benar
import { RouterLink } from 'vue-router'

const auth = useAuthStore()
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const clientError = ref('') // Untuk error validasi sisi klien

async function handleRegister() {
  clientError.value = ''; // Reset client error
  auth.error = null; // Reset auth store error

  if (password.value !== confirmPassword.value) {
    clientError.value = 'Passwords do not match.'
    return
  }
  if (password.value.length < 6) {
    clientError.value = 'Password should be at least 6 characters.'
    return;
  }

  await auth.signUp({
    email: email.value,
    password: password.value,
  })

  // Navigasi atau pesan sukses/error akan diurus di dalam action signUp di store
  // Jika signUp berhasil dan perlu konfirmasi, alert akan muncul dari store.
  // Jika langsung login, akan diarahkan ke Home.
  // Jika gagal, auth.error akan terisi.
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white p-10 rounded-xl shadow-lg">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create your account
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email-address" class="block text-sm font-medium text-gray-700 mb-1">Email address</label>
            <input
                id="email-address"
                name="email"
                type="email"
                v-model="email"
                autocomplete="email"
                required
                class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="user@company.com"
            />
          </div>
          <div class="pt-4">
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input
                id="password"
                name="password"
                type="password"
                v-model="password"
                autocomplete="new-password"
                required
                minlength="6"
                class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Password (min. 6 characters)"
            />
          </div>
          <div class="pt-4">
            <label for="confirm-password" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
            <input
                id="confirm-password"
                name="confirm-password"
                type="password"
                v-model="confirmPassword"
                autocomplete="new-password"
                required
                class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Confirm Password"
            />
          </div>
        </div>

        <div v-if="auth.error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
          <span class="block sm:inline">{{ auth.error }}</span>
        </div>

        <div v-if="clientError" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative" role="alert">
          <span class="block sm:inline">{{ clientError }}</span>
        </div>


        <div>
          <button
              type="submit"
              :disabled="auth.loading"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-700 disabled:opacity-50"
          >
            <span v-if="auth.loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 mr-3 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            Sign up
          </button>
        </div>
      </form>

      <p class="mt-8 text-center text-sm text-gray-600">
        Already have an account?
        <RouterLink :to="{ name: 'Login' }" class="font-medium text-indigo-600 hover:text-indigo-500">
          Log in
        </RouterLink>
      </p>
    </div>
  </div>
</template>

<style scoped>

</style>