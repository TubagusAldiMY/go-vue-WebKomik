// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
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
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL), // Gunakan createWebHistory untuk URL bersih
    routes,
})

export default router