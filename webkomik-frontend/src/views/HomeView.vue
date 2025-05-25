<script setup>
import { onMounted } from 'vue';
import { useComicStore } from '@/stores/comicStore';
import { useAuthStore } from '@/stores/authStore'; // Untuk tombol admin
import { RouterLink, useRouter } from 'vue-router'; // Import useRouter

const comicStore = useComicStore();
const auth = useAuthStore(); // Untuk cek peran admin
const router = useRouter(); // Untuk navigasi programatik

onMounted(() => {
  comicStore.fetchAllComics();
});

function goToCreateComicPage() {
  router.push({ name: 'AdminCreateComic' }); // <-- UBAH BARIS INI
}
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Telusuri Komik</h1>
      <button
          v-if="auth.isAuthenticated && auth.userRole === 'admin'"
          @click="goToCreateComicPage"
          class="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-2 px-4 rounded-lg shadow-md transition-colors duration-150"
      >
        + Tambah Komik Baru
      </button>
    </div>

    <div v-if="comicStore.loading" class="text-center py-10">
      <p class="text-lg text-gray-600">Memuat komik...</p>
    </div>

    <div v-else-if="comicStore.error" class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-md shadow-md" role="alert">
      <p class="font-bold">Oops! Terjadi kesalahan:</p>
      <p>{{ comicStore.error }}</p>
    </div>

    <div v-else-if="comicStore.comicsList.length === 0" class="text-center py-10">
      <p class="text-lg text-gray-600">Belum ada komik yang tersedia saat ini.</p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
      <div
          v-for="comic in comicStore.comicsList"
          :key="comic.id"
          class="bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300 ease-in-out"
      >
        <RouterLink :to="{ name: 'ComicDetail', params: { id: comic.id } }" class="block">
          <img
              :src="comic.cover_image_url || 'https://via.placeholder.com/300x450.png?text=No+Cover'"
              :alt="`Sampul ${comic.title}`"
              class="w-full h-72 object-cover"
          />
          <div class="p-4">
            <h3 class="text-lg font-semibold text-gray-800 truncate" :title="comic.title">{{ comic.title }}</h3>
            <p v-if="comic.author_name" class="text-sm text-gray-600 truncate" :title="comic.author_name">
              {{ comic.author_name }}
            </p>
            <p v-if="comic.genre_name" class="text-xs text-indigo-600 mt-1 bg-indigo-100 px-2 py-0.5 rounded-full inline-block">
              {{ comic.genre_name }}
            </p>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Style tambahan jika diperlukan */
.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>