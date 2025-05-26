// src/services/apiService.js
import { useAuthStore } from '@/stores/authStore'; // Untuk mendapatkan token jika diperlukan

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'; // Default ke port backend Go kita

async function request(endpoint, options = {}) {
    const authStore = useAuthStore(); // Akses store di dalam fungsi jika perlu token
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    // Tambahkan token otentikasi jika pengguna sudah login dan endpoint memerlukannya
    // Untuk saat ini, endpoint komik kita publik, jadi ini belum terlalu krusial
    // tapi ini adalah pola yang baik untuk endpoint yang dilindungi.
    if (authStore.session && authStore.session.access_token) {
        // Periksa apakah endpoint ini ada dalam daftar yang memerlukan auth atau tidak
        // Untuk contoh ini, kita asumsikan endpoint di bawah '/api/authRequired' atau serupa memerlukannya
        // Atau, Anda bisa menambahkan flag 'requiresAuth: true' ke options
        if (options.requiresAuth) {
            headers['Authorization'] = `Bearer ${authStore.session.access_token}`;
        }
    }

    const config = {
        ...options,
        headers,
    };

    try {
        console.log(`[apiService] Requesting URL: ${BASE_URL}${endpoint}`); // <--- TAMBAHKAN INI UNTUK DEBUGGING
        const response = await fetch(`${BASE_URL}${endpoint}`, config);
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: response.statusText }));
            // Buat error yang lebih informatif
            const error = new Error(errorData.message || `HTTP error! status: ${response.status}`);
            error.response = response; // Lampirkan seluruh respons ke objek error
            error.data = errorData;     // Lampirkan data error JSON ke objek error
            throw error;
        }
        if (response.status === 204) { // No Content
            return null;
        }
        return response.json();
    } catch (error) {
        console.error(`API request error for ${endpoint}:`, error);
        throw error; // Lempar ulang error agar bisa ditangani oleh pemanggil
    }
}

// Fungsi spesifik untuk endpoint komik
export const getComics = () => {
    return request('/comics'); // Endpoint GET /api/comics
};

export const getComicDetail = (id) => {
    return request(`/comics/${id}`); // Endpoint GET /api/comics/:id
};

export const createComic = (comicData) => {
    // Endpoint ini memerlukan otentikasi dan peran admin atau creator
    // Kita akan menandainya agar token ditambahkan
    return request('/comics', {
        method: 'POST',
        body: JSON.stringify(comicData),
        requiresAuth: true, // Tandai bahwa endpoint ini butuh token
    });
};

export const updateComic = (id, comicData) => {
    // Endpoint ini memerlukan otentikasi dan peran admin atau creator
    return request(`/comics/${id}`, {
        method: 'PUT',
        body: JSON.stringify(comicData),
        requiresAuth: true, // Tandai bahwa endpoint ini butuh token
    });
};

// Tambahkan fungsi lain di sini (deleteComic, dll.)