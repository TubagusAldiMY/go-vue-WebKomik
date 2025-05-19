// src/stores/comicStore.js
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getComics, getComicDetail } from '@/services/apiService'; // Pastikan path ini benar

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

    return {
        comicsList,
        currentComic,
        loading,
        error,
        fetchAllComics,
        fetchComicById,
        // createNewComic,
    };
});