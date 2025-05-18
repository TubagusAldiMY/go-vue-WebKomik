// src/main.js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router' // Pastikan path ini benar ke src/router/index.js
import './style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(router) // Router digunakan di sini
app.use(pinia)

app.mount('#app')
