<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Tambah Komik Baru</h1>

    <div class="max-w-2xl mx-auto bg-white p-8 rounded-lg shadow-xl">
      <form @submit.prevent="handleCreateComic">
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
          <label for="genre_id" class="block text-sm font-medium text-gray-700 mb-1">Genre ID (Contoh: 1)</label>
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

        <div v-if="comicStore.error" class="mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <span class="block sm:inline">{{ comicStore.error }}</span>
        </div>


        <div class="mt-6">
          <button
            type="submit"
            :disabled="comicStore.loading"
            class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline disabled:opacity-50 flex items-center justify-center"
          >
            <svg v-if="comicStore.loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ comicStore.loading ? 'Menyimpan...' : 'Simpan Komik' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useComicStore } from '@/stores/comicStore';
import { useRouter } from 'vue-router';

const comicStore = useComicStore();
const router = useRouter();

const comicForm = ref({
  title: '',
  description: null,
  author_name: null,
  genre_id: null,
  cover_image_url: null,
});

async function handleCreateComic() {
  comicStore.error = null; // Clear previous errors
  try {
    // Filter out null values or convert empty strings to null if backend expects that
    const payload = {};
    for (const key in comicForm.value) {
      if (comicForm.value[key] !== '' && comicForm.value[key] !== null) {
        payload[key] = comicForm.value[key];
      } else if (key === 'genre_id' && comicForm.value[key] === '') { 
         // Khusus untuk genre_id, jika string kosong, jangan kirim (atau kirim null eksplisit jika API-mu menanganinya)
         // Untuk input number, v-model.number akan menghasilkan 0 jika kosong, atau NaN jika tidak valid
         // Jadi, pastikan backend bisa menangani genre_id yang tidak ada atau null.
         // Untuk saat ini, jika kosong, tidak akan dikirim.
      } else {
         payload[key] = null; // Atau biarkan undefined jika backend lebih suka field tidak ada vs null
      }
    }
    if (payload.genre_id === 0) payload.genre_id = null; // Jika input number kosong, jadi 0, ubah ke null

    // Panggil action di store (akan kita buat/sempurnakan)
    // const newComic = await comicStore.createNewComic(payload);
    // Untuk sekarang, kita panggil langsung apiService via store (perlu action khusus di store)

    // Kita perlu action di comicStore untuk ini
    // Mari kita buat action createNewComic di comicStore.js
    const newComicData = await comicStore.createNewComic(payload); // Asumsi createNewComic ada di store

    if (newComicData) {
      alert('Komik berhasil ditambahkan!');
      router.push({ name: 'ComicDetail', params: { id: newComicData.id } }); // Arahkan ke detail komik baru
    }
    // Jika gagal, error akan di-set di store dan ditampilkan di template
  } catch (err) {
    // Error sudah di-set di store, tidak perlu tindakan tambahan di sini kecuali log
    console.error('Gagal di handleCreateComic:', err);
  }
}
</script>