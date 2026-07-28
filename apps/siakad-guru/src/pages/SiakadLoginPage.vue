<template>
  <q-page class="login-page bg-grey-2">
    <div class="login-wrapper row no-wrap">
      <!-- LEFT: BRANDING PANEL (Hidden on mobile) -->
      <div class="col-12 col-md-6 branding-panel column justify-between q-pa-xl">
        <div>
          <div class="row items-center q-gutter-sm">
            <q-avatar size="42px" color="secondary" text-color="white" icon="school" />
            <div class="text-h6 text-weight-bold text-white">SIAKAD</div>
          </div>
          <div class="text-caption text-grey-4 q-mt-xs letter-spacing-2">
            SISTEM INFORMASI AKADEMIK
          </div>
        </div>

        <div class="branding-content">
          <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
            SMK PASUNDAN JATINANGOR
          </div>
          <div class="text-h4 text-weight-bold text-white q-mt-sm branding-title">
            Satu Sistem, Seluruh Ekosistem Sekolah
          </div>
          <div class="text-body1 text-grey-4 q-mt-md branding-subtitle">
            Mengelola akademik, absensi, kesiswaan, dan kurikulum dalam satu platform terintegrasi.
          </div>
        </div>

        <div class="text-caption text-grey-5">
          Made with <q-icon name="favorite" color="red-10" size="xs" class="q-mx-xs" /> by. Deden
          Moh Jaenudin. © {{ currentYear }} Hak cipta dilindungi.
        </div>
      </div>

      <!-- RIGHT: LOGIN FORM -->
      <div class="col-12 col-md-6 form-panel column items-center justify-center q-pa-md">
        <q-card flat class="login-card">
          <q-card-section>
            <div class="text-h5 text-weight-bold text-primary">Masuk ke SIAKAD</div>
            <div class="text-body2 text-grey-7 q-mt-xs">
              Gunakan Email yang terdaftar di Dapodik.
            </div>
          </q-card-section>

          <q-card-section>
            <q-banner v-if="errorMessage" dense rounded class="bg-red-1 text-red-9 q-mb-md">
              <template v-slot:avatar>
                <q-icon name="error" color="red" />
              </template>
              {{ errorMessage }}
            </q-banner>

            <q-form @submit.prevent="handleLogin" class="q-gutter-md">
              <q-input v-model="form.identifier" outlined label="Email" lazy-rules :rules="[val => !!val || 'Wajib diisi',
              val => /.+@.+\..+/.test(val) || 'Format email tidak valid']" :disable="loading">
                <template v-slot:prepend>
                  <q-icon name="email" />
                </template>
              </q-input>

              <q-input v-model="form.password" outlined label="Kata Sandi" :type="showPassword ? 'text' : 'password'"
                lazy-rules :rules="[val => !!val || 'Wajib diisi']" :disable="loading" @keyup.enter="handleLogin">
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
                <template v-slot:append>
                  <q-icon :name="showPassword ? 'visibility_off' : 'visibility'" class="cursor-pointer"
                    @click="showPassword = !showPassword" />
                </template>
              </q-input>

              <div class="row items-center justify-between">
                <q-checkbox v-model="form.rememberMe" label="Ingat saya" dense :disable="loading" />
                <q-btn flat dense no-caps color="primary" label="Lupa kata sandi?" @click="showResetDialog = true"
                  :disable="loading" />
              </div>
              <div>
                <q-btn unelevated no-caps color="primary" label="Masuk" type="submit"
                  class="full-width q-py-sm text-weight-bold" :loading="loading" />

              </div>
            </q-form>
            <!-- ✨ AWAL AREA LOGIN ALTERNATIF BELAJAR.ID ✨ -->
            <LoginBelajarId />
            <!-- ✨ AKHIR AREA LOGIN ALTERNATIF BELAJAR.ID ✨ -->
          </q-card-section>

          <q-separator />

          <q-card-section class="text-center">
            <div class="text-caption text-grey-7">
              Kesulitan masuk? Hubungi Operator SIAKAD / Tata Usaha.
            </div>
          </q-card-section>
        </q-card>
        <!-- ✨ FOOTER MOBILE: Ditambahkan class 'lt-md' dan margin atas -->
        <!-- (Hanya muncul di layar hp/tablet kecil dan ditaruh di paling bawah form-panel) -->
        <div class="text-caption text-grey-7 lt-md q-mt-xs text-center">
          <span class="text-bold">Made with <q-icon name="favorite" color="red-10" size="xs" class="q-mx-xs" /> by.
            Deden
            Moh
            Jaenudin.</span>
          <span class="block">
            © {{ currentYear }} Hak cipta dilindungi.</span>
        </div>
      </div>
    </div>

    <!-- DIALOG: RESET PASSWORD -->
    <q-dialog v-model="showResetDialog">
      <q-card class="reset-dialog">
        <q-card-section>
          <div class="text-subtitle1 text-weight-bold text-primary">Reset Kata Sandi</div>
          <div class="text-caption text-grey-7 q-mt-xs">
            Masukkan alamat email. Permintaan reset akan diteruskan ke Operator SIAKAD untuk direset manual.
          </div>
        </q-card-section>
        <q-separator />
        <q-card-section>
          <q-input v-model="resetIdentifier" outlined dense label="Email" :disable="loadingReset" />
        </q-card-section>
        <q-separator />
        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat no-caps color="grey-7" label="Batal" v-close-popup :disable="loadingReset" />
          <q-btn unelevated no-caps color="primary" label="Kirim Permintaan" @click="submitReset"
            :loading="loadingReset" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'

import { useMeta } from 'quasar'
import { useUtils } from '@/composables/helper/useUtils'
import LoginBelajarId from '@/components/LoginBelajarId.vue'
import { useAuth } from '@/composables/useAuth'

// Mengubah title halaman secara dinamis
useMeta({
  title: 'SIAKAD'
})
const $q = useQuasar()
const router = useRouter()
const auth = useAuth()
const { currentYear } = useUtils()


const form = ref({
  identifier: '',
  password: '',
  rememberMe: false,
})

const showPassword = ref(false)
const loading = ref(auth.loading)
const errorMessage = ref('')
// Akan dipindahkan




// --- LOGIN LOGIC ---
const handleLogin = async () => {
  errorMessage.value = ''

  // Validasi manual tambahan (meski lazy-rules sudah menangani UI)
  if (!form.value.identifier || !form.value.password) {
    errorMessage.value = 'Email dan kata sandi wajib diisi.'
    return
  }

  try {
    // Memanggil action khusus SIAKAD di Auth Store
    const success = await auth.login({
      username: form.value.identifier,
      password: form.value.password,
      remember_me: form.value.rememberMe
    })

    // console.log('SiakadLoginPage', success)
    if (success) {
      $q.notify({ type: 'positive', message: 'Login berhasil!', position: 'top' })
      // Redirect
      router.push({ name: 'siakad-dashboard' })
    }
  } catch (error) {
    errorMessage.value = error.message || 'Terjadi kesalahan pada server.'
  }
}

// --- RESET PASSWORD LOGIC ---
const showResetDialog = ref(false)
const resetIdentifier = ref('')
const loadingReset = ref(false)

const submitReset = async () => {
  if (!resetIdentifier.value) {
    $q.notify({ type: 'warning', message: 'Masukkan Username terlebih dahulu.', position: 'top' })
    return
  }

  loadingReset.value = true
  try {
    // TODO: Nanti diganti dengan API call ke Golang
    // await api.post('/siakad/auth/forgot-password', { identifier: resetIdentifier.value })

    // Simulasi delay
    await new Promise(resolve => setTimeout(resolve, 1000))

    showResetDialog.value = false
    $q.notify({
      type: 'positive',
      message: 'Permintaan reset telah dikirim ke Operator SIAKAD.',
      position: 'top'
    })
    resetIdentifier.value = ''
  } catch (error) {
    console.log(error)
    $q.notify({ type: 'negative', message: 'Gagal mengirim permintaan.', position: 'top' })
  } finally {
    loadingReset.value = false
  }
}

// Opsional: Jika user sudah login, langsung lempar ke dashboard
onMounted(() => {
  // if (auth.isAuthenticated) {
  //   router.push('/siakad/dashboard')
  // }
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
}

.login-wrapper {
  min-height: 100vh;
}

.letter-spacing-2 {
  letter-spacing: 0.12em;
}

/* LEFT PANEL */
.branding-panel {
  background: linear-gradient(160deg, #0a192f 0%, #123a63 60%, #0a192f 100%);
  min-height: 100vh;
  position: relative;
}

.branding-content {
  max-width: 460px;
}

.branding-title {
  line-height: 1.3;
}

.branding-subtitle {
  max-width: 420px;
  line-height: 1.6;
}

/* RIGHT PANEL */
.form-panel {
  background-color: #f7f8fa;
  min-height: 100vh;
}

.login-card {
  width: 100%;
  max-width: 420px;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(10, 25, 47, 0.08);
}

.reset-dialog {
  min-width: 340px;
  max-width: 420px;
  width: 90vw;
}

/* Responsive: hide branding panel on small screens */
@media (max-width: 1023px) {
  .branding-panel {
    display: none;
  }

  .form-panel {
    width: 100%;
  }
}
</style>
