import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { supabase } from './lib/supabaseClient' // <-- Import supabase client

console.log('Supabase instance:', supabase) // <-- Tes logging

// Contoh sederhana mengambil data sesi (jika ada)
async function checkUserSession() {
    const { data: { session } } = await supabase.auth.getSession()
    if (session) {
        console.log('Sesi pengguna ditemukan:', session)
    } else {
        console.log('Tidak ada sesi pengguna aktif.')
    }
}
checkUserSession();

createApp(App).mount('#app')
