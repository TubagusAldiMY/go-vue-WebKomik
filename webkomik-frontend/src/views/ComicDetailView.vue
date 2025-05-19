<script setup>
import { computed, onMounted, ref } from 'vue';
import { useComicStore } from '@/stores/comicStore';
import { useRoute } from 'vue-router'; // Untuk mendapatkan parameter route

const props = defineProps({
  id: { // Menerima 'id' sebagai prop dari router
    type: [String, Number],
    required: true,
  }
});

const comicStore = useComicStore();
const route = useRoute(); // Alternatif untuk mendapatkan params jika props tidak digunakan

// Menggunakan computed property untuk mengakses currentComic dengan aman
const comic = computed(() => comicStore.currentComic);

const selectedChapterPages = ref([]);
const viewingChapterNumber = ref(null);

onMounted(() => {
  // Pastikan ID adalah angka sebelum memanggil store
  const comicId = Number(props.id || route.params.id);
  if (!isNaN(comicId)) {
    comicStore.fetchComicById(comicId);
  } else {
    comicStore.error = 'ID Komik tidak valid.';
  }
});

function viewChapter(comicId, chapterId) {
  const selectedComic = comicStore.currentComic;
  if (selectedComic && selectedComic.chapters) {
    const chapter = selectedComic.chapters.find(ch => ch.id === chapterId);
    if (chapter && chapter.pages) {
      selectedChapterPages.value = chapter.pages;
      viewingChapterNumber.value = chapter.chapter_number;
    } else {
      selectedChapterPages.value = [];
      viewingChapterNumber.value = null;
      console.warn(`Chapter dengan ID ${chapterId} atau halamannya tidak ditemukan.`);
    }
  }
}
</script>

<template>
  <div v-if="comicStore.loading" class="text-center py-10">
    <p class="text-lg text-gray-600">Memuat detail komik...</p>
  </div>
  <div v-else-if="comicStore.error" class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4" role="alert">
    <p class="font-bold">Error:</p>
    <p>{{ comicStore.error }}</p>
  </div>
  <div v-else-if="comic" class="container mx-auto px-4 py-8">
    <div class="bg-white shadow-xl rounded-lg overflow-hidden md:flex">
      <div class="md:w-1/3 p-4">
        <img
            :src="comic.cover_image_url || 'https://via.placeholder.com/400x600.png?text=No+Cover'"
            :alt="`Sampul ${comic.title}`"
            class="w-full h-auto object-contain rounded-md shadow"
        >
      </div>
      <div class="md:w-2/3 p-6">
        <h1 class="text-3xl md:text-4xl font-bold text-gray-800 mb-3">{{ comic.title }}</h1>
        <p v-if="comic.author_name" class="text-md text-gray-600 mb-1"><strong>Penulis:</strong> {{ comic.author_name }}</p>
        <p v-if="comic.genre_name" class="text-md text-gray-600 mb-4">
          <strong>Genre:</strong>
          <span class="text-indigo-600 bg-indigo-100 px-2 py-0.5 rounded-full text-sm">{{ comic.genre_name }}</span>
        </p>
        <p class="text-gray-700 leading-relaxed mb-6">{{ comic.description || 'Tidak ada deskripsi.' }}</p>

        <div class="border-t border-gray-200 pt-4">
          <h2 class="text-2xl font-semibold text-gray-700 mb-3">Chapter</h2>
          <div v-if="comic.chapters && comic.chapters.length > 0">
            <ul>
              <li
                  v-for="chapter in comic.chapters"
                  :key="chapter.id"
                  class="mb-2 p-3 bg-gray-50 hover:bg-gray-100 rounded-md transition-colors duration-150"
              >
                <a @click.prevent="viewChapter(comic.id, chapter.id)" href="#" class="block text-indigo-700 hover:underline">
                  Chapter {{ chapter.chapter_number }}: {{ chapter.title || 'Tanpa Judul' }}
                </a>
              </li>
            </ul>
          </div>
          <p v-else class="text-gray-500">Belum ada chapter untuk komik ini.</p>
        </div>
      </div>
    </div>

    <div v-if="selectedChapterPages.length > 0" class="mt-8 bg-white shadow-lg rounded-lg p-6">
      <h3 class="text-xl font-semibold mb-4">Halaman Chapter {{ viewingChapterNumber }}</h3>
      <div class="flex flex-col items-center">
        <img v-for="page in selectedChapterPages" :key="page.id" :src="page.image_url" :alt="`Halaman ${page.page_number}`" class="max-w-full md:max-w-3xl mb-2 border border-gray-300">
      </div>
    </div>

  </div>
  <div v-else class="text-center py-10">
    <p class="text-lg text-gray-600">Komik tidak ditemukan.</p>
  </div>
</template>

<style scoped>

</style>