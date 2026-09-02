<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Pengaturan Sistem</div>

    <q-tabs v-model="tab" dense class="q-mb-md" align="left">
      <q-tab name="umum" label="Umum" />
      <q-tab name="ujian" label="Pengaturan Ujian" />
      <q-tab name="email" label="Konfigurasi Email" />
      <q-tab name="tema" label="Tema & Branding" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <!-- Tab Umum -->
      <q-tab-panel name="umum" class="q-pa-none">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">Pengaturan Umum Aplikasi</div>
          </q-card-section>
          <q-card-section>
            <q-form @submit="saveUmum" class="q-gutter-md">
              <q-input v-model="settings.umum.nama_aplikasi" label="Nama Aplikasi" outlined
                :rules="[val => !!val || 'Nama aplikasi harus diisi']" />
              <q-input v-model="settings.umum.nama_institusi" label="Nama Institusi (Sekolah/Universitas)" outlined />
              <q-input v-model="settings.umum.tahun_ajaran_aktif" label="Tahun Ajaran Aktif" outlined
                hint="Contoh: 2025/2026" />
              <q-input v-model="settings.umum.semester_aktif" label="Semester Aktif" outlined hint="Ganjil / Genap" />
              <q-file v-model="settings.umum.logo" label="Logo Institusi" accept="image/*" outlined
                @update:model-value="handleLogoUpload">
                <template v-slot:prepend>
                  <q-icon name="image" />
                </template>
              </q-file>
              <div v-if="settings.umum.logoPreview" class="q-mt-sm">
                <img :src="settings.umum.logoPreview" style="max-height: 80px;" />
              </div>
              <div class="q-gutter-sm">
                <q-btn label="Simpan Pengaturan" type="submit" color="primary" />
                <q-btn label="Reset" flat color="negative" @click="resetUmum" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- Tab Pengaturan Ujian -->
      <q-tab-panel name="ujian" class="q-pa-none">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">Pengaturan Default Ujian</div>
          </q-card-section>
          <q-card-section>
            <q-form @submit="saveUjian" class="q-gutter-md">
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-6">
                  <q-input v-model.number="settings.ujian.durasi_default" label="Durasi Default (menit)" type="number"
                    outlined :rules="[val => val > 0 || 'Harus > 0']" />
                </div>
                <div class="col-12 col-md-6">
                  <q-select v-model="settings.ujian.tampilkan_hasil"
                    :options="['langsung', 'setelah_selesai', 'dijadwalkan']" label="Waktu Tampil Hasil" outlined
                    :option-label="opt => ({ langsung: 'Langsung setelah submit', setelah_selesai: 'Setelah ujian selesai semua', dijadwalkan: 'Dijadwalkan manual' }[opt])" />
                </div>
                <div class="col-12 col-md-6">
                  <q-toggle v-model="settings.ujian.bisa_ulang" label="Siswa dapat mengulang ujian" />
                </div>
                <div class="col-12 col-md-6">
                  <q-input v-model.number="settings.ujian.jumlah_ulang" label="Jumlah maksimal pengulangan"
                    type="number" outlined :disable="!settings.ujian.bisa_ulang" />
                </div>
                <div class="col-12 col-md-6">
                  <q-toggle v-model="settings.ujian.acak_soal" label="Acak urutan soal per siswa" />
                </div>
                <div class="col-12 col-md-6">
                  <q-toggle v-model="settings.ujian.acak_pilihan" label="Acak urutan pilihan jawaban" />
                </div>
                <div class="col-12 col-md-6">
                  <q-input v-model.number="settings.ujian.nilai_minimal_lulus" label="Nilai Minimal Lulus (%)"
                    type="number" outlined :rules="[val => val >= 0 && val <= 100]" />
                </div>
                <div class="col-12 col-md-6">
                  <q-toggle v-model="settings.ujian.auto_lock" label="Kunci layar saat ujian (anti cheating)" />
                </div>
              </div>
              <div class="q-gutter-sm">
                <q-btn label="Simpan Pengaturan" type="submit" color="primary" />
                <q-btn label="Reset" flat color="negative" @click="resetUjian" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- Tab Konfigurasi Email -->
      <q-tab-panel name="email" class="q-pa-none">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">Pengaturan Server Email</div>
          </q-card-section>
          <q-card-section>
            <q-form @submit="saveEmail" class="q-gutter-md">
              <q-input v-model="settings.email.host" label="SMTP Host" outlined hint="Contoh: smtp.gmail.com" />
              <q-input v-model.number="settings.email.port" label="SMTP Port" type="number" outlined />
              <q-select v-model="settings.email.encryption" :options="['tls', 'ssl', 'none']" label="Enkripsi"
                outlined />
              <q-input v-model="settings.email.username" label="Username / Email" outlined />
              <q-input v-model="settings.email.password" label="Password" type="password" outlined />
              <q-input v-model="settings.email.dari_email" label="Email Pengirim (From)" outlined type="email" />
              <q-input v-model="settings.email.dari_nama" label="Nama Pengirim" outlined />
              <div class="q-gutter-sm">
                <q-btn label="Simpan Pengaturan" type="submit" color="primary" />
                <q-btn label="Test Kirim Email" color="info" @click="testEmail" />
                <q-btn label="Reset" flat color="negative" @click="resetEmail" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- Tab Tema & Branding -->
      <q-tab-panel name="tema" class="q-pa-none">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">Kustomisasi Tema & Branding</div>
          </q-card-section>
          <q-card-section>
            <q-form @submit="saveTema" class="q-gutter-md">
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-6">
                  <q-input v-model="settings.tema.warna_primary" label="Warna Primary" type="color" outlined />
                  <div class="q-mt-sm"
                    :style="{ backgroundColor: settings.tema.warna_primary, width: '100%', height: '30px', borderRadius: '4px' }">
                  </div>
                </div>
                <div class="col-12 col-md-6">
                  <q-input v-model="settings.tema.warna_secondary" label="Warna Secondary" type="color" outlined />
                  <div class="q-mt-sm"
                    :style="{ backgroundColor: settings.tema.warna_secondary, width: '100%', height: '30px', borderRadius: '4px' }">
                  </div>
                </div>
                <div class="col-12 col-md-6">
                  <q-input v-model="settings.tema.warna_accent" label="Warna Accent" type="color" outlined />
                  <div class="q-mt-sm"
                    :style="{ backgroundColor: settings.tema.warna_accent, width: '100%', height: '30px', borderRadius: '4px' }">
                  </div>
                </div>
                <div class="col-12 col-md-6">
                  <q-toggle v-model="settings.tema.dark_mode" label="Mode Gelap (Dark Mode)" />
                </div>
                <div class="col-12 col-md-6">
                  <q-select v-model="settings.tema.font"
                    :options="['Roboto', 'Poppins', 'Open Sans', 'Lato', 'Montserrat']" label="Font Utama" outlined />
                </div>
                <div class="col-12 col-md-6">
                  <q-input v-model="settings.tema.logo_header" label="URL Logo Header" outlined
                    hint="URL gambar untuk ditampilkan di header" />
                </div>
              </div>
              <div class="q-gutter-sm">
                <q-btn label="Simpan Tema" type="submit" color="primary" />
                <q-btn label="Reset ke Default" flat color="negative" @click="resetTema" />
              </div>
            </q-form>
          </q-card-section>
          <q-card-section>
            <div class="text-subtitle1">Preview Tema</div>
            <div class="row q-col-gutter-md q-mt-sm">
              <div class="col-12">
                <q-card
                  :style="{ backgroundColor: settings.tema.warna_primary + '20', borderTop: `4px solid ${settings.tema.warna_primary}` }">
                  <q-card-section>
                    <div class="text-h6" :style="{ color: settings.tema.warna_primary }">Contoh Judul</div>
                    <div :style="{ color: settings.tema.warna_secondary }">Warna sekunder untuk teks</div>
                    <q-btn :style="{ backgroundColor: settings.tema.warna_primary, color: '#fff' }"
                      label="Tombol Primary" class="q-mt-sm" />
                    <q-btn :style="{ backgroundColor: settings.tema.warna_secondary, color: '#fff' }"
                      label="Tombol Secondary" class="q-mt-sm q-ml-sm" />
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// Default settings
const defaultSettings = {
  umum: {
    nama_aplikasi: 'CBT Ujian Online',
    nama_institusi: 'SMA Negeri 1',
    tahun_ajaran_aktif: '2025/2026',
    semester_aktif: 'Ganjil',
    logo: null,
    logoPreview: null
  },
  ujian: {
    durasi_default: 90,
    tampilkan_hasil: 'setelah_selesai',
    bisa_ulang: false,
    jumlah_ulang: 1,
    acak_soal: true,
    acak_pilihan: false,
    nilai_minimal_lulus: 70,
    auto_lock: false
  },
  email: {
    host: 'smtp.gmail.com',
    port: 587,
    encryption: 'tls',
    username: '',
    password: '',
    dari_email: '',
    dari_nama: 'CBT System'
  },
  tema: {
    warna_primary: '#1976D2',
    warna_secondary: '#26A69A',
    warna_accent: '#9C27B0',
    dark_mode: false,
    font: 'Roboto',
    logo_header: ''
  }
}

const settings = reactive({
  umum: { ...defaultSettings.umum },
  ujian: { ...defaultSettings.ujian },
  email: { ...defaultSettings.email },
  tema: { ...defaultSettings.tema }
})

const tab = ref('umum')

// Load from localStorage
function loadSettings() {
  const saved = localStorage.getItem('cbt_settings')
  if (saved) {
    const parsed = JSON.parse(saved)
    Object.assign(settings.umum, parsed.umum || defaultSettings.umum)
    Object.assign(settings.ujian, parsed.ujian || defaultSettings.ujian)
    Object.assign(settings.email, parsed.email || defaultSettings.email)
    Object.assign(settings.tema, parsed.tema || defaultSettings.tema)
  }
  // Terapkan tema ke root (opsional, untuk demo)
  applyTheme()
}

function saveAllSettings() {
  localStorage.setItem('cbt_settings', JSON.stringify({
    umum: settings.umum,
    ujian: settings.ujian,
    email: settings.email,
    tema: settings.tema
  }))
  applyTheme()
}

function applyTheme() {
  // Untuk demo, kita ubah variabel CSS root
  document.documentElement.style.setProperty('--q-primary', settings.tema.warna_primary)
  document.documentElement.style.setProperty('--q-secondary', settings.tema.warna_secondary)
  document.documentElement.style.setProperty('--q-accent', settings.tema.warna_accent)
  if (settings.tema.dark_mode) {
    document.body.classList.add('body--dark')
  } else {
    document.body.classList.remove('body--dark')
  }
  // Font
  document.body.style.fontFamily = settings.tema.font
}

// Tab Umum
function handleLogoUpload(file) {
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      settings.umum.logoPreview = e.target.result
      settings.umum.logo = file
    }
    reader.readAsDataURL(file)
  } else {
    settings.umum.logoPreview = null
    settings.umum.logo = null
  }
}
function saveUmum() {
  saveAllSettings()
  $q.notify({ type: 'positive', message: 'Pengaturan umum disimpan' })
}
function resetUmum() {
  settings.umum = { ...defaultSettings.umum }
  saveAllSettings()
  $q.notify({ type: 'info', message: 'Pengaturan umum direset' })
}

// Tab Ujian
function saveUjian() {
  saveAllSettings()
  $q.notify({ type: 'positive', message: 'Pengaturan ujian disimpan' })
}
function resetUjian() {
  settings.ujian = { ...defaultSettings.ujian }
  saveAllSettings()
  $q.notify({ type: 'info', message: 'Pengaturan ujian direset' })
}

// Tab Email
function saveEmail() {
  saveAllSettings()
  $q.notify({ type: 'positive', message: 'Konfigurasi email disimpan' })
}
function resetEmail() {
  settings.email = { ...defaultSettings.email }
  saveAllSettings()
  $q.notify({ type: 'info', message: 'Konfigurasi email direset' })
}
function testEmail() {
  // Simulasi test email
  $q.loading.show({ message: 'Mengirim email test...' })
  setTimeout(() => {
    $q.loading.hide()
    $q.notify({ type: 'info', message: 'Email test dikirim ke ' + settings.email.dari_email + ' (simulasi)' })
  }, 1500)
}

// Tab Tema
function saveTema() {
  saveAllSettings()
  $q.notify({ type: 'positive', message: 'Tema disimpan' })
}
function resetTema() {
  settings.tema = { ...defaultSettings.tema }
  saveAllSettings()
  $q.notify({ type: 'info', message: 'Tema direset ke default' })
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
/* Optional: style untuk preview warna */
</style>
