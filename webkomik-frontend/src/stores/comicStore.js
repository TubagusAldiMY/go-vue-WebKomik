// src/stores/comicStore.js
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getComics, getComicDetail, createComic as apiCreateComic } from '@/services/apiService'; // Pastikan path ini benar

export const useComicStore = defineStore('comics', () => {
    // State
    const comicsList = ref([]);
    const currentComic = ref(null); // Untuk menyimpan detail satu komik
    const loading = ref(false);
    const error = ref(null);

    // Actions
    async function fetchAllComics() {
        loading.value = true;
        error.value = null;
        comicsList.value = []; // Kosongkan dulu untuk menghindari data lama jika ada error
        try {
            const response = await getComics(); // response sudah dalam format { data: [...] }
            comicsList.value = response.data || []; // Ambil array dari field 'data'
        } catch (e) {
            error.value = e.message || 'Gagal mengambil daftar komik.';
            console.error('Error fetching comics:', e);
        } finally {
            loading.value = false;
        }
    }

    async function fetchComicById(id) {
        loading.value = true;
        error.value = null;
        currentComic.value = null; // Kosongkan dulu
        try {
            const response = await getComicDetail(id); // response sudah dalam format { data: { ... } }
            currentComic.value = response.data || null; // Ambil objek dari field 'data'
        } catch (e) {
            if (e.response && e.response.status === 404) {
                error.value = 'Komik tidak ditemukan.';
            } else {
                error.value = e.message || `Gagal mengambil detail komik ID: ${id}.`;
            }
            console.error(`Error fetching comic ID ${id}:`, e);
        } finally {
            loading.value = false;
        }
    }

    // Action untuk membuat komik (membutuhkan otentikasi admin)
    // Ini akan dipanggil dari halaman/komponen yang sesuai
    // async function createNewComic(comicData) {
    //   loading.value = true;
    //   error.value = null;
    //   try {
    //     const response = await createComic(comicData); // createComic dari apiService
    //     // Mungkin tambahkan komik baru ke comicsList atau fetch ulang semua
    //     await fetchAllComics(); // Contoh: fetch ulang daftar
    //     return response.data; // Kembalikan data komik yang baru dibuat
    //   } catch (e) {
    //     error.value = e.data?.details || e.message || 'Gagal membuat komik baru.';
    //     console.error('Error creating comic:', e);
    //     throw e; // Lempar ulang error agar bisa ditangani di komponen
    //   } finally {
    //     loading.value = false;
    //   }
    // }

    async function createNewComic(comicData) {
    loading.value = true;
    error.value = null;
    try {
      // Pastikan data yang dikirim hanya yang diperlukan oleh CreateComicInput di backend
      // atau model Comic di backend.
      const payload = {
        title: comicData.title,
        description: comicData.description || null,
        author_name: comicData.author_name || null,
        genre_id: comicData.genre_id || null,
        cover_image_url: comicData.cover_image_url || null,
      };

      const response = await apiCreateComic(payload); // Panggil fungsi dari apiService

      // Setelah berhasil, mungkin kita ingin refresh daftar komik atau navigasi
      await fetchAllComics(); // Contoh: fetch ulang daftar komik di halaman Beranda
      return response.data; // Kembalikan data komik yang baru dibuat (dari field 'data' di respons)
    } catch (e) {
      // apiService sudah melempar error yang lebih baik
      // e.data mungkin berisi detail error validasi dari backend
      if (e.data && e.data.details) {
        error.value = `Input tidak valid: ${e.data.details}`;
      } else if (e.data && e.data.error) {
        error.value = e.data.error;
      }
      else {
        error.value = e.message || 'Gagal membuat komik baru.';
      }
      console.error('Error creating comic in store:', e);
      throw e; // Lempar ulang error agar bisa ditangani lebih lanjut jika perlu
    } finally {
      loading.value = false;
    }
  }
  
  async function updateExistingComic(id, comicData) {
    loading.value = true;
    error.value = null;
    try {
      // Hanya kirim field yang memang ada nilainya (untuk update parsial)
      const payload = {};
      if (comicData.title !== undefined && comicData.title !== null) payload.title = comicData.title;
      if (comicData.description !== undefined) payload.description = comicData.description;
      if (comicData.author_name !== undefined) payload.author_name = comicData.author_name;
      if (comicData.genre_id !== undefined) payload.genre_id = comicData.genre_id || null;
      if (comicData.cover_image_url !== undefined) payload.cover_image_url = comicData.cover_image_url;
      
      // Panggil fungsi updateComic dari apiService
      const response = await updateComic(id, payload);
      
      // Perbarui daftar komik dan data komik saat ini jika sama
      await fetchAllComics();
      if (currentComic.value && currentComic.value.id === id) {
        currentComic.value = response.data;
      }
      
      return response.data; // Kembalikan data komik yang telah diperbarui
    } catch (e) {
      // Handle error seperti pada createNewComic
      if (e.data && e.data.details) {
        error.value = `Input tidak valid: ${e.data.details}`;
      } else if (e.data && e.data.error) {
        error.value = e.data.error;
      } else {
        error.value = e.message || `Gagal memperbarui komik ID ${id}.`;
      }
      console.error(`Error updating comic ID ${id}:`, e);
      throw e;
    } finally {
      loading.value = false;
    }
  }

  return {
    comicsList,
    currentComic,
    loading,
    error,
    fetchAllComics,
    fetchComicById,
    createNewComic,
    updateExistingComic, // Expose new update action
  };
});
