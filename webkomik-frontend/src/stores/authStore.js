// src/stores/authStore.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { supabase } from '../lib/supabaseClient' // Pastikan path ini benar
import router from '../router' // Import router untuk navigasi

export const useAuthStore = defineStore('auth', () => {
    // State
    const user = ref(null) // Akan menyimpan data user dari Supabase (profile, dll.)
    const session = ref(null) // Akan menyimpan data sesi dari Supabase
    const loading = ref(false)
    const error = ref(null)
    const returnUrl = ref(null) // Untuk redirect setelah login
    const userRole = ref(null)

    // Getters (computed properties)
    const isAuthenticated = computed(() => !!session.value && !!session.value.user)
    const userEmail = computed(() => session.value?.user?.email || null)
    const userId = computed(() => session.value?.user?.id || null)
    const isAdmin = computed(() => userRole.value === 'admin');


     // Fungsi untuk parse token dan ekstrak role (jika ada)
  // Supabase JWT biasanya memiliki custom claims di dalam user.app_metadata atau user.user_metadata
  // Atau jika Anda mengaturnya sebagai root claim.
  // Untuk Supabase, peran biasanya ada di `session.value.user.app_metadata.role` atau `session.value.user.user_metadata.role`
  // atau jika Anda telah membuat fungsi SQL untuk menambahkannya sebagai top-level claim.
  // Kita asumsikan role ada di `app_metadata.role` atau `user_metadata.role` untuk contoh ini.
    function updateUserRoleFromSession() {
    if (session.value && session.value.user) {
        // Coba dari app_metadata dulu
        if (session.value.user.app_metadata && session.value.user.app_metadata.role) {
        userRole.value = session.value.user.app_metadata.role;
        }
        // Jika tidak ada, coba dari user_metadata
        else if (session.value.user.user_metadata && session.value.user.user_metadata.role) {
        userRole.value = session.value.user.user_metadata.role;
        }
        else {
        userRole.value = 'user'; // Default
        }
    } else {
        userRole.value = null;
    }
    }
    // Actions
    async function fetchUserSession() {
        loading.value = true
        error.value = null
        try {
            const { data, error: sessionError } = await supabase.auth.getSession()
            if (sessionError) throw sessionError
            session.value = data.session
            if (data.session?.user) {
                user.value = data.session.user;
                updateUserRoleFromSession();
            } else {
                user.value = null
                userRole.value = null // Reset userRole jika tidak ada sesi
            }
        } catch (e) {
            error.value = e.message
            console.error('Error fetching session:', e.message)
        } finally {
            loading.value = false
        }
    }

    async function loginWithPassword(credentials) {
        loading.value = true
        error.value = null
        try {
            const { data, error: loginError } = await supabase.auth.signInWithPassword({
                email: credentials.email,
                password: credentials.password,
            })
            if (loginError) throw loginError
            if (data.session && data.user) {
                session.value = data.session;
                user.value = data.user;
                updateUserRoleFromSession(); // <-- PANGGIL FUNGSI INI
    }

            // Navigasi setelah login berhasil
            if (returnUrl.value) {
                router.push(returnUrl.value)
                returnUrl.value = null // Reset returnUrl
            } else {
                router.push({ name: 'Home' }) // Redirect ke halaman Home atau dashboard
            }

        } catch (e) {
            error.value = e.message
            console.error('Login error:', e.message)
            // Pastikan user dan session di-null-kan jika login gagal
            session.value = null
            user.value = null
            userRole.value = null;
        } finally {
            loading.value = false
        }
    }

    async function signUp(credentials) {
        loading.value = true;
        error.value = null;
        try {
            const { data, error: signUpError } = await supabase.auth.signUp({
                email: credentials.email,
                password: credentials.password,
                // Anda bisa menambahkan options seperti data untuk user_metadata di sini jika perlu
                // options: {
                //   data: {
                //     full_name: credentials.fullName, // contoh
                //   }
                // }
            });
            if (signUpError) throw signUpError;

            // Supabase mungkin mengirim email konfirmasi.
            // Sesi dan user mungkin tidak langsung ada sampai email dikonfirmasi,
            // tergantung pengaturan Supabase Anda (Secure email change & Email confirmations).
            session.value = data.session;
            user.value = data.user;

            // Beri pesan atau arahkan sesuai kebutuhan setelah sign up
            // Misalnya, jika perlu konfirmasi email:
            if (data.user && !data.session && data.user.identities && data.user.identities.length > 0) {
                alert('Pendaftaran berhasil! Silakan cek email Anda untuk konfirmasi.');
                router.push({ name: 'Login' }); // Arahkan ke login setelah info
            } else if (data.session) {
                // Jika sesi langsung ada (misal konfirmasi email dimatikan)
                router.push({ name: 'Home' });
            } else {
                // Skenario lain, mungkin perlu penanganan khusus
                alert('Pendaftaran diproses. Ikuti instruksi selanjutnya jika ada.');
                router.push({ name: 'Login' });
            }

        } catch (e) {
            error.value = e.message;
            console.error('Sign up error:', e.message);
        } finally {
            loading.value = false;
        }
    }


    async function logout() {
        loading.value = true
        error.value = null
        try {
            const { error: logoutError } = await supabase.auth.signOut()
            if (logoutError) throw logoutError
            user.value = null
            session.value = null
            userRole.value = null;
            router.push({ name: 'Login' }) // Redirect ke halaman Login setelah logout
        } catch (e) {
            error.value = e.message
            console.error('Logout error:', e.message)
        } finally {
            loading.value = false
        }
    }

    // Panggil fetchUserSession saat store diinisialisasi untuk memeriksa sesi yang ada
    // fetchUserSession() // Panggil ini di App.vue atau saat mounting utama

    return {
        user,
        session,
        loading,
        error,
        isAuthenticated,
        userEmail,
        userId,
        userRole,
        isAdmin,
        returnUrl,
        fetchUserSession,
        loginWithPassword,
        signUp,
        logout,
    }
})