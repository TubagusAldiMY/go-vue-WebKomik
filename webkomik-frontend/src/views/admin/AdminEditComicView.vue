<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Edit Komik: {{ comicForm.title }}</h1>

    <div v-if="loading" class="flex justify-center items-center min-h-[300px]">
      <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else-if="!currentComic && !error" class="flex justify-center items-center min-h-[300px]">
      <p class="text-gray-600">Memuat data komik...</p>
    </div>

    <div v-else-if="error" class="mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <span class="block sm:inline">{{ error }}</span>
    </div>

    <div v-else class="max-w-2xl mx-auto bg-white p-8 rounded-lg shadow-xl">
      <form @submit.prevent="handleUpdateComic">
        <div class="mb-4">
          <label for="title" class="block text-sm font-medium text-gray-700 mb-1">Judul Komik</label>
          <input
            type="text"
            id="title"
            v-model="comicForm.title"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <div class="mb-4">
          <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
          <textarea
            id="description"
            v-model="comicForm.description"
            rows="4"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          ></textarea>
        </div>

        <div class="mb-4">
          <label for="author_name" class="block text-sm font-medium text-gray-700 mb-1">Nama Penulis</label>
          <input
            type="text"
            id="author_name"
            v-model="comicForm.author_name"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <div class="mb-4">
          <label for="genre_id" class="block text-sm font-medium text-gray-700 mb-1">Genre ID</label>
          <input
            type="number"
            id="genre_id"
            v-model.number="comicForm.genre_id"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <div class="mb-6">
          <label for="cover_image_url" class="block text-sm font-medium text-gray-700 mb-1">URL Gambar Sampul</label>
          <input
            type="url"
            id="cover_image_url"
            v-model="comicForm.cover_image_url"
            placeholder="https://example.com/image.jpg"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <div class="flex items-center space-x-4 mt-6">
          <button
            type="submit"
            :disabled="comicStore.loading"
            class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline disabled:opacity-50 flex items-center justify-center"
          >
            <svg v-if="comicStore.loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ comicStore.loading ? 'Menyimpan...' : 'Simpan Perubahan' }}
          </button>
          <button
            type="button"
            @click="router.back()"
            class="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline"
          >
            Batal
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useComicStore } from '@/stores/comicStore';
import { useRouter, useRoute } from 'vue-router';

const comicStore = useComicStore();
const router = useRouter();
const route = useRoute();

// Dapatkan ID komik dari URL
const comicId = parseInt(route.params.id);

// Reactive data
const comicForm = ref({
  title: '',
  description: '',
  author_name: '',
  genre_id: null,
  cover_image_url: '',
});

// Computed properties untuk memudahkan akses ke state dari store
const loading = computed(() => comicStore.loading);
const error = computed(() => comicStore.error);
const currentComic = computed(() => comicStore.currentComic);

// Ambil data komik dan isi form saat komponen dimuat
onMounted(async () => {
  try {
    await comicStore.fetchComicById(comicId);
    // Isi formulir dengan data komik yang sudah ada
    if (comicStore.currentComic) {
      comicForm.value.title = comicStore.currentComic.title || '';
      comicForm.value.description = comicStore.currentComic.description || '';
      comicForm.value.author_name = comicStore.currentComic.author_name || '';
      comicForm.value.genre_id = comicStore.currentComic.genre_id || null;
      comicForm.value.cover_image_url = comicStore.currentComic.cover_image_url || '';
    }
  } catch (err) {
    console.error('Gagal mengambil detail komik:', err);
  }
});

// Fungsi untuk mengirim perubahan ke server
async function handleUpdateComic() {
  comicStore.error = null; // Reset error sebelumnya
  
  try {
    // Siapkan payload yang hanya berisi field yang diubah
    const payload = {};
    
    // Bandingkan dengan data asli untuk menentukan field yang diubah
    const original = comicStore.currentComic;
    
    // Hanya tambahkan field yang berubah ke payload
    if (original.title !== comicForm.value.title) {
      payload.title = comicForm.value.title;
    }
    
    // Untuk field opsional, perlu penanganan khusus
    if (original.description !== comicForm.value.description) {
      payload.description = comicForm.value.description || null;
    }
    
    if (original.author_name !== comicForm.value.author_name) {
      payload.author_name = comicForm.value.author_name || null;
    }
    
    if (original.genre_id !== comicForm.value.genre_id) {
      payload.genre_id = comicForm.value.genre_id || null;
    }
    
    if (original.cover_image_url !== comicForm.value.cover_image_url) {
      payload.cover_image_url = comicForm.value.cover_image_url || null;
    }
    
    // Panggil action store untuk update komik
    const updatedComic = await comicStore.updateExistingComic(comicId, payload);
    
    if (updatedComic) {
      alert('Komik berhasil diperbarui!');
      router.push({ name: 'ComicDetail', params: { id: comicId } });
    }
  } catch (err) {
    // Error sudah ditangani di store
    console.error('Gagal di handleUpdateComic:', err);
  }
}
</script>
