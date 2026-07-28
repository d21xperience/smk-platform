<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-mx-auto q-px-md max-width-center">
      <!-- 1. HEADER HALAMAN -->
      <div class="text-center q-mb-xl">
        <div class="text-h4 text-sm-h3 text-weight-bold text-primary">
          Formulir Pendaftaran <span class="text-amber">PPDB</span>
        </div>
        <p class="text-body1 text-grey-7 q-mt-sm">
          Tahun Ajaran 2026/2027 — Silakan isi data calon siswa secara lengkap dan benar.
        </p>
      </div>

      <!-- 2. WORKSPACE FORMULIR MULTI-STEP (QSTEPPER) -->
      <div class="row justify-center">
        <div class="col-12 col-md-9">
          <q-card flat bordered class="bg-white rounded-borders shadow-2">
            <q-stepper
              v-model="step"
              ref="stepper"
              color="primary"
              animated
              header-nav
              :vertical="$q.screen.xs"
              class="no-shadow"
            >
              <!-- STEP 1: PILIHAN JURUSAN & DATA DIRI -->
              <q-step :name="1" title="Data Calon Siswa" icon="person" :done="step > 1">
                <q-form ref="formStep1" class="q-gutter-md q-pt-md">
                  <div class="row q-col-gutter-md">
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >Pilihan Komp. Keahlian (Jurusan)</q-item-label
                      >
                      <q-select
                        outlined
                        dense
                        v-model="studentData.major"
                        :options="majorOptions"
                        emit-value
                        map-options
                        :rules="[(val) => !!val || 'Jurusan wajib dipilih']"
                      />
                    </div>
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >NISN (Nomor Induk Siswa Nasional)</q-item-label
                      >
                      <q-input
                        outlined
                        dense
                        v-model="studentData.nisn"
                        mask="##########"
                        placeholder="10 Digit NISN Anda"
                        :rules="[
                          (val) => (val && val.length === 10) || 'NISN wajib 10 digit angka',
                        ]"
                      />
                    </div>
                  </div>

                  <div>
                    <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                      >Nama Lengkap Siswa (Sesuai Ijazah)</q-item-label
                    >
                    <q-input
                      outlined
                      dense
                      v-model="studentData.fullName"
                      placeholder="Contoh: ADITYA PUTRA"
                      :rules="[(val) => !!val || 'Nama lengkap wajib diisi']"
                    />
                  </div>

                  <div class="row q-col-gutter-md">
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >Tempat Lahir</q-item-label
                      >
                      <q-input
                        outlined
                        dense
                        v-model="studentData.pob"
                        placeholder="Contoh: Bandung"
                        :rules="[(val) => !!val || 'Tempat lahir wajib diisi']"
                      />
                    </div>
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >Tanggal Lahir</q-item-label
                      >
                      <q-input outlined dense v-model="studentData.dob" placeholder="YYYY-MM-DD">
                        <template v-slot:append>
                          <q-icon name="event" class="cursor-pointer">
                            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                              <q-date v-model="studentData.dob" mask="YYYY-MM-DD">
                                <div class="row items-center justify-end">
                                  <q-btn v-close-popup label="Selesai" color="primary" flat />
                                </div>
                              </q-date>
                            </q-popup-proxy>
                          </q-icon>
                        </template>
                      </q-input>
                    </div>
                  </div>

                  <div>
                    <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                      >Sekolah Asal (SMP/MTs)</q-item-label
                    >
                    <q-input
                      outlined
                      dense
                      v-model="studentData.prevSchool"
                      placeholder="Contoh: SMP Negeri 1 Majalaya"
                      :rules="[(val) => !!val || 'Sekolah asal wajib diisi']"
                    />
                  </div>
                </q-form>
              </q-step>

              <!-- STEP 2: DATA ORANG TUA / WALI -->
              <q-step
                :name="2"
                title="Data Orang Tua / Wali"
                icon="family_restroom"
                :done="step > 2"
              >
                <q-form ref="formStep2" class="q-gutter-md q-pt-md">
                  <div>
                    <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                      >Nama Ayah / Ibu Kandung / Wali</q-item-label
                    >
                    <q-input
                      outlined
                      dense
                      v-model="parentData.parentName"
                      placeholder="Contoh: BUDI SANTOSO"
                      :rules="[(val) => !!val || 'Nama orang tua wajib diisi']"
                    />
                  </div>

                  <div class="row q-col-gutter-md">
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >No. WhatsApp Aktif (Penting)</q-item-label
                      >
                      <q-input
                        outlined
                        dense
                        v-model="parentData.phone"
                        placeholder="Contoh: 08123456789"
                        :rules="[
                          (val) => !!val || 'Nomor WhatsApp wajib diisi untuk info kelulusan',
                        ]"
                      />
                    </div>
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >Pekerjaan Orang Tua / Wali</q-item-label
                      >
                      <q-input
                        outlined
                        dense
                        v-model="parentData.occupation"
                        placeholder="Contoh: Karyawan Swasta"
                      />
                    </div>
                  </div>

                  <div>
                    <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                      >Alamat Rumah Lengkap</q-item-label
                    >
                    <q-input
                      outlined
                      dense
                      type="textarea"
                      v-model="parentData.address"
                      placeholder="Tuliskan nama jalan, RT/RW, kelurahan, dan kecamatan"
                      :rules="[(val) => !!val || 'Alamat lengkap wajib diisi']"
                    />
                  </div>
                </q-form>
              </q-step>

              <!-- STEP 3: UPLOAD DOKUMEN BERKAS -->
              <q-step :name="3" title="Unggah Berkas Pendukung" icon="cloud_upload">
                <div class="q-pt-md">
                  <q-banner
                    class="bg-amber-1 text-grey-9 rounded-borders q-mb-lg text-caption"
                    inline-actions
                  >
                    <template v-slot:avatar><q-icon name="info" color="amber-9" /></template>
                    Format berkas wajib berupa Gambar (JPG/PNG) atau PDF dengan ukuran maksimal
                    masing-masing berkas 2MB.
                  </q-banner>

                  <div class="row q-col-gutter-md">
                    <!-- Upload Kartu Keluarga -->
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >Upload Scan Kartu Keluarga (KK)</q-item-label
                      >
                      <q-uploader
                        url="https://httpbin.org"
                        label="Klik atau Tarik Berkas KK"
                        accept=".jpg, .jpeg, .png, .pdf"
                        max-file-size="2097152"
                        auto-upload
                        flat
                        bordered
                        class="full-width"
                      />
                    </div>
                    <!-- Upload Ijazah / SKL -->
                    <div class="col-12 col-sm-6">
                      <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                        >Upload Scan Ijazah / SKL (Surat Keterangan Lulus)</q-item-label
                      >
                      <q-uploader
                        url="https://httpbin.org"
                        label="Klik atau Tarik Berkas Ijazah"
                        accept=".jpg, .jpeg, .png, .pdf"
                        max-file-size="2097152"
                        auto-upload
                        flat
                        bordered
                        class="full-width"
                      />
                    </div>
                  </div>
                </div>
              </q-step>

              <!-- TOMBOL NAVIGASI DI BAGIAN BAWAH KONTROL STEPPER -->
              <template v-slot:navigation>
                <q-separator class="q-mt-xl q-mb-md" />
                <div class="row justify-between items-center">
                  <!-- Tombol Kembali -->
                  <q-btn
                    v-if="step > 1"
                    flat
                    color="primary"
                    @click="$refs.stepper.previous()"
                    label="Langkah Sebelumnya"
                    class="q-px-md text-weight-bold"
                  />
                  <div v-else></div>

                  <!-- Tombol Lanjut / Selesai -->
                  <q-btn
                    @click="nextStep"
                    color="primary"
                    :label="step === 3 ? 'Kirim Pendaftaran' : 'Langkah Berikutnya'"
                    icon-right="navigate_next"
                    class="q-px-lg text-weight-bold"
                    :loading="isSubmitting"
                    unelevated
                  />
                </div>
              </template>
            </q-stepper>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()

const step = ref(1)
const isSubmitting = ref(false)

// Objek Penampung Data Formulir PPDB
const studentData = ref({
  major: 'tbsm',
  nisn: '',
  fullName: '',
  pob: '',
  dob: '2010-01-01',
  prevSchool: '',
})
const parentData = ref({ parentName: '', phone: '', occupation: '', address: '' })

// Pilihan Opsi Dropdown 5 Jurusan Sinkron
const majorOptions = [
  { label: 'Teknik & Bisnis Sepeda Motor (TBSM)', value: 'tbsm' },
  { label: 'Teknik Komputer & Jaringan (TKJ)', value: 'tkj' },
  { label: 'Teknik Kendaraan Ringan (TKR)', value: 'tkr' },
  { label: 'Akuntansi & Keuangan Lembaga (AKL)', value: 'akl' },
  { label: 'Otomatisasi & Tata Kelola Perkantoran (OTKP)', value: 'otkp' },
]
// Logika Kendali Validasi Sebelum Pindah Step / Submit Akhir
function nextStep() {
  if (step.value === 1) {
    // Validasi data diri calon siswa terlebih dahulu sebelum lanjut step 2
    step.value = 2
  } else if (step.value === 2) {
    // Validasi data orang tua sebelum lanjut step 3 upload berkas
    step.value = 3
  } else if (step.value === 3) {
    // Eksekusi Submit Formulir Akhir (Kirim Data)submit
    submitRegistration()
  }
}
function submitRegistration() {
  isSubmitting.value = true
  // Simulasi pemrosesan pengiriman data ke server API selama 2 detik
  setTimeout(() => {
    isSubmitting.value = false // Berikan pop-up alert sukses yang interaktif
    $q.dialog({
      title: 'Pendaftaran Berhasil!',
      message:
        'Data PPDB online Anda telah sukses tersimpan di sistem kami. Kode pendaftaran resmi dan bukti tanda terima telah dikirimkan otomatis ke nomor WhatsApp Orang Tua Anda.',
      ok: { label: 'Kembali Ke Beranda', color: 'primary', unelevated: true },
      persistent: true,
    }).onOk(() => {
      // Tendang balik pengguna ke halaman depan utama
      router.push('/')
    })
  }, 2000)
}
</script>
<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.no-shadow {
  box-shadow: none !important;
}
</style>
