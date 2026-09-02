<template>
  <q-layout view="lHh Lpr lFf" class="text-white overflow-md-hidden">
    <!-- Top Navbar Mini sesuai gambar -->
    <!-- <q-header class="bg-blue-10 q-py-xs floating-header absolute fixed-top">
      <q-toolbar class="row justify-between items-center q-px-md">
        <div class="row items-center q-gutter-sm">
          <q-avatar size="32px" square>
            <img src="https://wikimedia.org" alt="Logo Kemen">
          </q-avatar>
          <div class="text-caption text-weight-bold text-amber-4">Aplikasi e-Rapor SMK</div>
        </div>
        <div class="row q-gutter-sm text-caption text-weight-medium items-center">
          <q-btn flat round dense :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'"
            :color="$q.dark.isActive ? 'amber-4' : 'white'" @click="toggleDarkMode">
            <q-tooltip class="bg-grey-9 text-white">
              {{ $q.dark.isActive ? 'Mode Terang' : 'Mode Gelap' }}
            </q-tooltip>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header> -->

    <!-- Konten Utama: Split Layar Kiri & Kanan -->
    <q-page-container>
      <q-page class="row items-stretch justify-center">

        <!-- SISI KIRI: Banner Informasi & Ilustrasi (60% Lebar pada Desktop) -->
        <div class="col-12 col-md-7 flex flex-center bg-grey-2 text-blue-10 relative-position q-pa-sm">
          <!-- Section 1: Hero Banner -->
          <section id="hero" class="row items-center justify-center q-py-xl q-px-md bg-gradient">
            <div class="col-12 col-md-5 q-pa-md hidden-md hidden-xs">
              <q-icon name="dashboard_customize" size="100px" color="primary" class="opacity-50" />
            </div>
            <div class="col-12  q-pa-md text-center text-md-left">
              <q-badge color="black" text-color="white" label="Tahun Ajaran 2026/2027"
                class="q-py-xs q-px-sm q-mb-md" />
              <h1 class="text-h3 text-md-h2 text-weight-bolder q-my-none leading-tight">
                Sistem Informasi Akademik Modern
              </h1>
              <p class="text-subtitle1 text-grey-7 q-mt-md q-mb-xl">
                Integrasi data administrasi, nilai, kehadiran, dan pembelajaran dalam satu platform digital terpadu
                untuk efisiensi sekolah masa kini.
              </p>
            </div>

          </section>
        </div>

        <!-- SISI KANAN: Form Login Box Putih (40% Lebar) -->
        <div class="col-12 col-md-5 bg-grey-2 flex flex-center q-pa-lg">
          <q-card flat class="login-box shadow-2 q-pa-xl bg-white text-grey-9"
            style="width: 100%; max-width: 480px; border-radius: 8px;">
            <div class="text-center q-mb-lg">
              <div class="text-subtitle1 text-weight-bold text-grey-8 row items-center justify-center q-gutter-xs">
                <q-icon name="login" color="blue-10" />
                <span>MASUK</span>
              </div>
              <div class="text-caption text-grey-15">Silahkan Masuk Untuk Memulai Aplikasi</div>
              <q-separator class="q-my-md" />
            </div>

            <!-- Form -->
            <q-form @submit.prevent="onLoginSubmit" class="q-gutter-md">
              <!-- Username -->
              <div class="row items-center">
                <div class="col-3 text-caption text-weight-bold text-grey-7">Username</div>
                <div class="col-9">
                  <q-input v-model="form.username" outlined dense placeholder="Isi dengan pengguna"
                    :rules="[val => !!val || 'Required']" hide-bottom-space />
                </div>
              </div>

              <!-- Password -->
              <div class="row items-center">
                <div class="col-3 text-caption text-weight-bold text-grey-7">Password</div>
                <div class="col-9">
                  <q-input v-model="form.password" :type="showPwd ? 'text' : 'password'" outlined dense
                    placeholder="Ketikkan Kata Sandi" :rules="[val => !!val || 'Required']" hide-bottom-space>
                    <template v-slot:append>
                      <q-icon :name="showPwd ? 'visibility_off' : 'visibility'" class="cursor-pointer"
                        @click="showPwd = !showPwd" />
                    </template>
                  </q-input>
                </div>
              </div>
              <!-- Lupa Password -->
              <div class="row justify-end">
                <router-link to="/forgot" class="text-caption text-blue-10 text-weight-bold decoration-none">Lupa
                  Password
                  ?</router-link>
              </div>

              <!-- Tombol Masuk -->
              <div class="q-mt-lg">
                <q-btn type="submit" color="blue-10" icon="lock" label="Masuk" class="full-width text-weight-bold"
                  unevaluated :loading="submitting" />
              </div>
            </q-form>

            <div class="text-center text-caption text-grey-5 q-mt-md">
              Isi User & Password dengan benar
            </div>
          </q-card>
        </div>

      </q-page>
    </q-page-container>


    <!-- Footer -->
    <q-footer class="text-white q-py-md bg-blue-10">
      <div class="text-center text-caption">
        <div class="q-mt-xs">SIAKAD SMK &copy; 2026. (Made with <span class="text-red-14">❤</span> by. Deden M.J.)</div>
      </div>
    </q-footer>
  </q-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { mockAuthService } from '@/services/mocks/authMock'

const $q = useQuasar()

const showPwd = ref(false)
const loadingMetadata = ref(true)
const submitting = ref(false)

const schoolOptions = ref([])
const semesterOptions = ref([])

const form = reactive({
  username: '',
  password: '',
  school: null,
  semester: null
})

// Mengambil metadata (Sekolah & Semester) saat halaman dimuat
onMounted(async () => {
  try {
    const [resSchools, resSemesters] = await Promise.all([
      mockAuthService.getSchools(),
      mockAuthService.getSemesters()
    ])
    schoolOptions.value = resSchools.data
    semesterOptions.value = resSemesters.data

    // Set default value jika data tersedia sesuai gambar
    if (schoolOptions.value.length > 0) form.school = schoolOptions.value[0].value
    if (semesterOptions.value.length > 0) form.semester = semesterOptions.value[0].value
  } catch (err) {
    console.error('Gagal memuat metadata', err)
  } finally {
    loadingMetadata.value = false
  }
})

async function onLoginSubmit() {
  submitting.value = true
  try {
    const res = await mockAuthService.login(form)
    $q.notify({
      type: 'positive',
      message: 'Login Berhasil',
      caption: `Selamat datang di ${res.data.user.name}`,
      position: 'top-right'
    })
    // Lakukan router.push('/dashboard') ke dashboard utama Anda disini
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: 'Login Gagal',
      caption: error.response?.data?.message || 'Koneksi bermasalah',
      position: 'top-right'
    })
  } finally {
    submitting.value = false
  }
}
</script>
<style scoped>
/* ==========================================================================
   MOBILE-FIRST DEFAULT STYLES (Layar < 992px)
   ========================================================================== */
.main-page-row {
  min-height: calc(100vh - 50px);
  /* Fleksibel untuk mobile agar bisa di-scroll */
}

.left-banner-panel {
  display: none;
  /* Mobile-first: Sembunyikan banner kiri */
}

.full-width-mobile {
  width: 100%;
}

.text-none {
  text-decoration: none;
}

.line-height-tight {
  line-height: 1.2;
}

/* ==========================================================================
   DESKTOP OPTIMIZATION STYLES (Layar >= 992px / MD ke Atas) - BEBAS SCROLL
   ========================================================================== */
@media (min-width: 992px) {

  /* 1. Kunci total tinggi layout induk agar tidak melar melebihi layar monitor */
  :deep(.q-layout) {
    height: 100vh !important;
    max-height: 100vh !important;
    overflow: hidden !important;
  }

  /* 2. Jadikan posisi header absolute agar lepas dari aliran dokumen (tidak mendorong komponen bawah) */
  .floating-header {
    position: absolute !important;
    top: 0;
    left: 0;
    right: 0;
    z-index: 2000;
    background: transparent !important;
    box-shadow: none !important;
  }

  /* 3. Paksa container mengabaikan padding-top otomatis hasil kalkulasi skrip Quasar */
  :deep(.custom-container) {
    padding-top: 0px !important;
    padding-bottom: 0px !important;
    height: 100vh !important;
    max-height: 100vh !important;
    overflow: hidden !important;
  }

  /* 4. Bagi porsi tinggi row halaman tepat mengisi sisa ruang di atas footer */
  .main-page-row {
    height: calc(100vh - 32px) !important;
    max-height: calc(100vh - 32px) !important;
    overflow: hidden !important;
  }

  /* 5. Tampilkan kembali banner kiri dengan padding pengaman atas agar teks tidak tertabrak menu */
  .left-banner-panel {
    display: flex !important;
    padding-top: 90px !important;
  }

  .branding-box {
    top: 50px;
    left: 40px;
  }

  /* 6. Kunci posisi footer tetap di paling bawah */
  .footer-panel {
    position: fixed !important;
    bottom: 0;
    left: 0;
    right: 0;
    height: 32px;
    z-index: 2000;
  }
}

.max-width-content {
  max-width: 460px;
}

.border-top {
  border-top: 1px solid rgba(255, 255, 255, 0.15);
}

.opacity-70 {
  opacity: 0.7;
}
</style>
