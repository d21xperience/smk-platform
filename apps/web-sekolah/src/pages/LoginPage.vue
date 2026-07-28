<template>
  <q-layout view="lHh Lpr lFf">
    <q-page class="login-page flex flex-center">
      <!-- KOTAK UTAMA LOGIN -->
      <q-card flat bordered class="login-card overflow-hidden">
        <div class="login-card__inner">

          <!-- PANEL BRANDING -->
          <!-- Mobile: banner horizontal ramping di atas. Desktop (>=1024px): panel vertikal di sisi kiri -->
          <div class="brand-panel bg-primary text-white relative-position">
            <div class="brand-panel__content z-top">
              <!-- <q-avatar size="56px" class="bg-white q-pa-xs shadow-2 brand-panel__logo"> -->
              <img src="~@/assets/logo-smk.png" alt="Logo Sekolah" onerror="this.style.display='none'" width="90px" />
              <!-- </q-avatar> -->
              <div class="brand-panel__text">
                <div class="text-weight-bold text-amber brand-panel__title">Portal Admin</div>
                <div class="text-caption text-grey-3 brand-panel__subtitle">
                  Sistem Manajemen Konten - SMK Pasundan Jatinangor
                </div>
              </div>
            </div>
            <div class="absolute-full bg-pattern"></div>
          </div>

          <!-- PANEL FORMULIR -->
          <div class="form-panel">
            <div class="text-h6 text-weight-bold text-grey-8">Selamat Datang Kembali</div>
            <div class="text-caption text-grey-6 q-mb-lg">Silakan masuk menggunakan akun resmi Anda.</div>

            <q-form @submit="handleLogin" class="q-gutter-md">
              <!-- Input Email / NIP -->
              <q-input outlined v-model="loginForm.identifier" label="Email" placeholder="contoh: budi@sekolah.sch.id"
                autocomplete="username" lazy-rules :rules="[val => val && val.length > 0 || 'Email wajib diisi']">
                <template v-slot:prepend>
                  <q-icon name="email" color="primary" />
                </template>
              </q-input>

              <!-- Input Password dengan Fitur Show/Hide -->
              <q-input outlined v-model="loginForm.password" :type="isPasswordVisible ? 'text' : 'password'"
                label="Kata Sandi" autocomplete="current-password" lazy-rules
                :rules="[val => val && val.length > 0 || 'Kata sandi wajib diisi']">
                <template v-slot:prepend>
                  <q-icon name="lock" color="primary" />
                </template>
                <template v-slot:append>
                  <q-btn flat round dense :icon="isPasswordVisible ? 'visibility_off' : 'visibility'"
                    :aria-label="isPasswordVisible ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'"
                    @click="isPasswordVisible = !isPasswordVisible" />
                </template>
              </q-input>


              <!-- TOMBOL SUBMIT -->
              <div class="q-pt-sm">
                <q-btn label="Masuk Ke Dashboard" type="submit" color="primary" no-caps
                  class="full-width text-weight-bold submit-btn" :loading="isLoading" unelevated />
              </div>
            </q-form>

            <!-- FOOTER LINK -->
            <div class="text-center q-mt-lg text-caption text-grey-6">
              Lupa kata sandi? Hubungi
              <span class="text-primary text-weight-bold cursor-pointer">Tim IT / Tata Usaha</span>
            </div>
          </div>

        </div>
      </q-card>
    </q-page>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()

const isLoading = ref(false)
const isPasswordVisible = ref(false)

// Struktur Form Input
const loginForm = ref({
  identifier: '',
  password: '',
  simulatedRole: 'guru' // Default role uji coba
})


// Eksekusi Logika Login
function handleLogin() {
  isLoading.value = true

  // TODO: ganti simulasi ini dengan pemanggilan API POST /api/auth/login ke Go backend
  setTimeout(() => {
    isLoading.value = false

    // 1. Simpan Informasi Session ke LocalStorage (Simulasi Token)
    localStorage.setItem('user_token', 'mock-jwt-token-xyz')
    localStorage.setItem('user_role', loginForm.value.simulatedRole)
    localStorage.setItem('user_name', loginForm.value.identifier.split('@')[0] || 'Staf Sekolah')

    // 2. Berikan Notifikasi Pop-up Sukses
    $q.notify({
      type: 'positive',
      message: `Login Berhasil! Selamat datang kembali, ${localStorage.getItem('user_name')}`,
      position: 'top',
      timeout: 2000
    })

    // 3. Alihkan Navigasi ke Ruang Kerja Admin Panel (Fase 2)
    router.push('/admin/dashboard')
  }, 1500)
}
</script>

<style scoped>
/* ============================================================
   MOBILE FIRST: gaya dasar di bawah ini menyasar layar kecil.
   Media query min-width menambah/menimpa gaya untuk layar besar.
   ============================================================ */

.login-page {
  min-height: 100vh;
  padding: 16px;
  padding-top: max(16px, env(safe-area-inset-top));
  padding-bottom: max(16px, env(safe-area-inset-bottom));
  box-sizing: border-box;
}

.login-card {
  width: 100%;
  max-width: 960px;
  border-radius: 16px;
  box-shadow: 0 2px 10px rgba(13, 43, 78, 0.12);
}

.login-card__inner {
  display: flex;
  flex-direction: column;
}

/* --- Panel Branding: banner ramping di mobile --- */
.brand-panel {
  padding: 20px 16px;
  background: linear-gradient(135deg, #0d2b4e 0%, #1a4a7a 100%);
}

.brand-panel__content {
  display: flex;
  align-items: center;
  gap: 14px;
  text-align: left;
}

.brand-panel__text {
  min-width: 0;
}

.brand-panel__title {
  font-size: 1.05rem;
  line-height: 1.2;
}

.brand-panel__subtitle {
  display: none;
  /* disembunyikan di layar sangat kecil agar banner tetap ringkas */
}

/* --- Panel Form --- */
.form-panel {
  padding: 24px 20px 20px;
}

.dev-mode-box {
  border: 1px dashed #d0d0d0;
  border-radius: 8px;
}

.submit-btn {
  min-height: 48px;
  /* target sentuh nyaman di mobile */
}

/* ============================================================
   >= 600px (tablet kecil ke atas): tampilkan subjudul branding
   ============================================================ */
@media (min-width: 600px) {
  .brand-panel__subtitle {
    display: block;
    margin-top: 2px;
  }

  .form-panel {
    padding: 32px 40px 28px;
  }
}

/* ============================================================
   >= 1024px (desktop): layout split kiri-kanan seperti aslinya
   ============================================================ */
@media (min-width: 1024px) {
  .login-page {
    padding: 24px;
  }

  .login-card {
    box-shadow: 0 8px 30px rgba(13, 43, 78, 0.18);
  }

  .login-card__inner {
    flex-direction: row;
    min-height: 520px;
  }

  .brand-panel {
    flex: 0 0 42%;
    max-width: 42%;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 48px 32px;
  }

  .brand-panel__content {
    flex-direction: column;
    text-align: center;
  }

  .brand-panel__logo {
    width: 80px !important;
    height: 80px !important;
  }

  .brand-panel__title {
    font-size: 1.5rem;
    margin-top: 16px;
  }

  .brand-panel__subtitle {
    margin-top: 8px;
  }

  .form-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 48px 56px;
  }
}

/* Dekorasi Latar Samping Kiri */
.bg-pattern {
  background-image: radial-gradient(circle at 20% 30%, #ffffff 1px, transparent 1px),
    radial-gradient(circle at 75% 70%, #ffffff 1px, transparent 1px);
  background-size: 20px 20px;
  opacity: 0.15;
  pointer-events: none;
}
</style>
