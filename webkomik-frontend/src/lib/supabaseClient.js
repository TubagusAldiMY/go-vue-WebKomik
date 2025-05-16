// src/lib/supabaseClient.js
import { createClient } from '@supabase/supabase-js'

// Ambil environment variables dari import.meta.env
const supabaseUrl = import.meta.env.VITE_SUPABASE_URL
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY

// Periksa apakah variabel ada (opsional tapi bagus untuk debugging)
if (!supabaseUrl || !supabaseAnonKey) {
    console.error(
        'Supabase URL atau Anon Key tidak ditemukan. Pastikan sudah diatur di file .env dan diawali dengan VITE_'
    )
}

// Buat dan ekspor Supabase client
export const supabase = createClient(supabaseUrl, supabaseAnonKey)