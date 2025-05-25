// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
// Kita akan membuat komponen view ini nanti
// import HomeView from '../views/HomeView.vue'
// import LoginView from '../views/LoginView.vue'
// import RegisterView from '../views/RegisterView.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        // component: HomeView // Akan kita aktifkan nanti
        component: () => import('../views/HomeView.vue') // Lazy load
    },
    {
        path: '/login',
        name: 'Login',
        // component: LoginView
        component: () => import('../views/LoginView.vue') // Lazy load
    },
    {
        path: '/register',
        name: 'Register',
        // component: RegisterView
        component: () => import('../views/RegisterView.vue') // Lazy load
    },
    // Tambahkan route lain di sini, misalnya untuk detail komik:
    {
        path: '/comic/:id', // :id adalah parameter dinamis
        name: 'ComicDetail',
        component: () => import('../views/ComicDetailView.vue'), // Lazy load
        props: true // Ini akan meneruskan route.params sebagai props ke komponen
    },
    // === ADMIN ROUTES ===
    {
        path: '/admin/comics/create',
        name: 'AdminCreateComic',
        component: () => import('../views/admin/AdminCreateComicView.vue'), // Buat view ini
        meta: { requiresAuth: true, requiresAdmin: true } // Meta untuk navigation guard
    },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation Guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Pastikan sesi sudah coba di-fetch sebelum guard berjalan
  // Ini penting jika pengguna me-refresh halaman di route yang dilindungi
  if (authStore.session === null && !authStore.loading) { // Cek juga loading agar tidak fetch berulang saat inisialisasi
    await authStore.fetchUserSession()
  }

  const isAuthenticated = authStore.isAuthenticated
  const isAdmin = authStore.userRole === 'admin' // Asumsi userRole sudah ada di authStore

  if (to.meta.requiresAuth && !isAuthenticated) {
    // Jika route memerlukan login dan pengguna belum login,
    // simpan URL yang dituju agar bisa redirect setelah login
    authStore.returnUrl = to.fullPath
    next({ name: 'Login' })
  } else if (to.meta.requiresAdmin && (!isAuthenticated || !isAdmin)) {
    // Jika route memerlukan admin dan pengguna bukan admin (atau belum login)
    // Bisa arahkan ke halaman "Unauthorized" atau Home
    console.warn('Akses ditolak: Membutuhkan peran admin.');
    next({ name: 'Home' }) // Atau buat halaman 'Unauthorized'
  } else if (to.meta.requiresGuest && isAuthenticated) {
    // Jika route hanya untuk tamu (belum login) tapi pengguna sudah login
    next({ name: 'Home' }) // Arahkan ke Home
  }
  else {
    next() // Lanjutkan navigasi
  }
})

export default router