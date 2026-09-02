<template>
  <q-layout view="lHh Lpr lFf">
    <q-page-container>
      <q-page class="flex flex-center bg-grey-2">
        <div class="login-container">
          <q-card class="login-card q-pa-sm" style="width: 100%; max-width: 450px">
            <q-card-section class="text-center">
              <img src="~assets/logo-smk.png" alt="Logo" style="height: 80px" />
              <div class="text-h5">SIAKAD</div>
              <!-- <div class="text-subtitle2 text-grey-6">Sistem Informasi Akademik Terpadu</div> -->
              <div class="text-subtitle2 text-grey-6">SMK Pasundan Jatinangor</div>
            </q-card-section>

            <q-card-section>
              <q-form @submit="handleLogin" class="">
                <q-input v-model="form.username" label="Username" outlined dense autofocus
                  :rules="[val => !!val || 'Username harus diisi']">
                  <template v-slot:prepend>
                    <q-icon name="person" />
                  </template>
                </q-input>

                <q-input v-model="form.password" label="Password" type="password" outlined dense
                  :rules="[val => !!val || 'Password harus diisi']" class="q-mt-lg q-mb-md">
                  <template v-slot:prepend>
                    <q-icon name="lock" />
                  </template>
                </q-input>

                <div class="row justify-between items-center">
                  <q-checkbox v-model="rememberMe" label="Ingat saya" dense />
                  <a href="#" class="text-caption text-primary" @click.prevent="lupaPassword">Lupa password?</a>
                </div>

                <q-btn type="submit" label="MASUK" color="primary" class="full-width q-mt-lg" :loading="loading"
                  rounded />
              </q-form>
              <q-separator class="q-my-md" />

              <div class="text-center text-grey-6 q-mb-md">
                atau masuk menggunakan
              </div>

              <q-btn outline color="primary" class="full-width"
                icon="img:https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg"
                label="Masuk dengan Belajar.id" @click="loginWithGoogle" />
            </q-card-section>

            <q-card-section class="text-center q-pt-none">
              <div class="text-caption text-grey-6">
                &copy; {{ new Date().getFullYear() }} SMK Pasundan. All rights reserved.
              </div>
            </q-card-section>
          </q-card>

          <!-- Pesan error -->
          <q-banner v-if="errorMessage" class="bg-negative text-white fixed-top"
            style="top: 20px; left: 50%; transform: translateX(-50%); width: auto; min-width: 300px; z-index: 1000">
            <template v-slot:avatar>
              <q-icon name="error" />
            </template>
            {{ errorMessage }}
          </q-banner>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useUserStore } from 'stores/user'
import { useAuthStore } from 'stores/auth-store'

const $q = useQuasar()
const router = useRouter()
const userStore = useUserStore()
const authStore = useAuthStore()

const form = ref({ username: '', password: '' })
const rememberMe = ref(false)
const loading = ref(false)
const errorMessage = ref('')

// Jika sudah login, redirect ke dashboard sesuai role
onMounted(() => {
  authStore.loadFromStorage()
  if (authStore.isLoggedIn) {
    redirectToDashboard(authStore.user.role)
  }
  userStore.loadData() // pastikan data user tersedia
})

function handleLogin() {
  loading.value = true
  errorMessage.value = ''

  // Simulasi delay (bisa diganti dengan panggilan API)
  setTimeout(() => {
    const user = userStore.getUserByUsername(form.value.username)
    if (user && user.password === form.value.password && user.is_active) {
      // Login berhasil
      authStore.login(user, 'dummy-token')
      if (rememberMe.value) {
        localStorage.setItem('remembered_user', form.value.username)
      } else {
        localStorage.removeItem('remembered_user')
      }
      redirectToDashboard(user.role)
    } else {
      errorMessage.value = 'Username atau password salah, atau akun tidak aktif.'
      loading.value = false
      setTimeout(() => { errorMessage.value = '' }, 3000)
    }
    loading.value = false
  }, 800)
}

function redirectToDashboard(role) {
  // Arahkan berdasarkan role
  const roleMap = {
    admin: '/admin/dashboard',
    operator: '/operator/dashboard',
    kepala_sekolah: '/kepsek/dashboard',
    wakasek: '/wakasek/dashboard',
    kaprodi: '/kaprodi/dashboard',
    guru: '/guru/dashboard',
    siswa: '/siswa/dashboard',
    orangtua: '/orangtua/dashboard'
  }
  const path = roleMap[role] || '/'
  router.push(path)
}

function lupaPassword() {
  $q.dialog({
    title: 'Lupa Password',
    message: 'Silakan hubungi administrator sekolah untuk reset password.',
    ok: 'OK'
  })
}

const loginWithGoogle = () => {
  console.log('hello world')
}
</script>

<style scoped>
.login-container {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.login-card {
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}
</style>
